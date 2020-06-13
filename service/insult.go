package service

import (
	"fmt"
	"github.com/juju/loggo"
	"insultService/model"
	"insultService/repository"
	"math/rand"
	"sort"
	"strings"
)

// Insult is an interface that contains methods relating to insults
type Insult interface {
	GenerateInsult(who model.Users, rank string) (message *string, id *string, err error)
	GetInsultsStats() (insultStat *model.InsultStat, err error)
	GetUserInfo(userID string) (userInfo *model.UserInfo, err error)
	IncreaseUserExperience(userID string) (message *string, err error)
	randomWordChooser(words *model.Words, adjCount int) (adjective, noun, verb string)
}

type insult struct {
	fireStore repository.FireStore
	log       loggo.Logger
	rand      rand.Rand
}

// NewInsult creates a new insult service
func NewInsult(fire repository.FireStore, log loggo.Logger, rand rand.Rand) Insult {
	return &insult{
		fireStore: fire,
		log:       log,
		rand:      rand,
	}
}

// GenerateInsult returns a string with a generated insult and an error bubbled up from firestore if any
func (i *insult) GenerateInsult(who model.Users, rank string) (message *string, id *string, err error) {
	//Should generate an Insult
	titles, err := i.fireStore.ReadTitles()
	if err != nil {
		return nil, nil, err
	}
	words, err := i.fireStore.ReadAllWords()
	if err != nil {
		return nil, nil, err
	}

	adjCount := 1
	i.log.Infof("%v", titles.Titles)
	for _, title := range titles.Titles {
		if rank != title.Name {
			adjCount++
		} else if rank == title.Name {
			break
		}
	}

	adj, noun, verb := i.randomWordChooser(words, adjCount)
	insultContents := model.InsultContent{
		Verb:      verb,
		Adjective: adj,
		Noun:      noun,
	}
	insult := insultMessage(who, adj, noun, verb, i.rand)
	//Should insert generated insult into firebase collection
	id, err = i.fireStore.InsertInsultEntry(insultContents)
	//Should produce an error if failed insert, but still return proper insult
	if err != nil {
		return &insult, nil, err
	}

	return &insult, id, nil
}

func (i *insult) IncreaseUserExperience(userID string) (*string, error) {
	//Check User Info and if they don't exist, save a new update
	titles, err := i.fireStore.ReadTitles()
	userInfoList, err := i.fireStore.ReadUserInfo([]model.QueryArg{{Path: "username", Op: "==", Value: userID}})
	if err != nil {
		i.log.Errorf("error retrieving profile with corresponding username and password")
		return nil, err
	}
	if len(userInfoList) == 0 {
		_, err = i.generateNewUser(userID, titles)
		if err != nil {
			return nil, err
		}
		message := fmt.Sprintf("Welcome!  Successfully Created a new user of the Insult Bot!  Your Current Rank is: %s", titles.Titles[0].Name)
		return &message, nil
	}
	//If user exists, increment their exp by 1
	userInfo := userInfoList[0]
	userInfo.Experience++
	//Check to see if user's current exp matches a new rank requirement
	newRank := determineTitle(userInfo.Experience, *titles)
	var message *string
	if newRank != userInfo.Rank {
		rankMessage := fmt.Sprintf("Congrats!  You have obtained the rank of: %s", newRank)
		userInfo.Rank = newRank
		message = &rankMessage
	}
	err = i.fireStore.UpdateUserInfo(userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to update user with error: %v", err)
	}
	return message, nil
}

func (i *insult) GetInsultsStats() (insultStat *model.InsultStat, err error) {
	insultContents, err := i.fireStore.ReadInsults()
	var verbs map[string]int = make(map[string]int)
	var nouns map[string]int = make(map[string]int)
	var adjectives map[string]int = make(map[string]int)
	for _, insultContent := range insultContents {
		verbs[insultContent.Verb]++
		nouns[insultContent.Noun]++
		adjectives[insultContent.Adjective]++
	}
	var verbArray []model.VerbCount
	for k, v := range verbs {
		verbArray = append(verbArray, model.VerbCount{
			Verb:  k,
			Count: v,
		})
	}
	var adjectiveArray []model.AdjectiveCount
	for k, v := range adjectives {
		adjectiveArray = append(adjectiveArray, model.AdjectiveCount{
			Adjective: k,
			Count:     v,
		})
	}
	var nounArray []model.NounCount
	for k, v := range nouns {
		nounArray = append(nounArray, model.NounCount{
			Noun:  k,
			Count: v,
		})
	}
	sort.Slice(verbArray, func(i, j int) bool {
		return verbArray[i].Count < verbArray[j].Count
	})
	sort.Slice(adjectiveArray, func(i, j int) bool {
		return adjectiveArray[i].Count < adjectiveArray[j].Count
	})
	sort.Slice(nounArray, func(i, j int) bool {
		return nounArray[i].Count < nounArray[j].Count
	})
	return &model.InsultStat{
		Adjectives: adjectiveArray,
		Verbs:      verbArray,
		Nouns:      nounArray,
	}, nil
}

func (i *insult) GetUserInfo(userID string) (userInfo *model.UserInfo, err error) {
	titles, err := i.fireStore.ReadTitles()
	userInfoList, err := i.fireStore.ReadUserInfo([]model.QueryArg{{Path: "username", Op: "==", Value: userID}})
	if err != nil {
		i.log.Errorf("error retrieving user with this username")
		return nil, err
	}
	if len(userInfoList) == 0 {
		userInfo, err = i.generateNewUser(userID, titles)
		if err != nil {
			return nil, err
		}
		return userInfo, nil
	}

	return userInfoList[0], nil
}

func (i *insult) generateNewUser(userID string, titles *model.Titles) (*model.UserInfo, error) {
	userInfo := model.UserInfo{
		Name:       userID,
		Rank:       determineTitle(0, *titles),
		Experience: 0,
	}
	i.log.Infof("userInfo: %v", userInfo)
	err := i.fireStore.UpdateUserInfo(&userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create user with error: %v", err)
	}
	return &userInfo, nil
}

func determineTitle(exp int, titles model.Titles) string {
	rank := ""
	for _, title := range titles.Titles {
		if exp >= title.Experience {
			rank = title.Name
		}
	}
	return rank
}

func (i *insult) randomWordChooser(words *model.Words, adjCount int) (adjective, noun, verb string) {
	adjList := ""
	for word := 0; word < adjCount; word++ {
		adjList += words.Adjective[i.rand.Intn(len(words.Adjective))]
		if word != adjCount-1 {
			adjList += ", "
		}
	}
	adjective = strings.TrimSpace(adjList)
	noun = words.Noun[i.rand.Intn(len(words.Noun))]
	verb = words.Verb[i.rand.Intn(len(words.Verb))]
	return adjective, noun, verb
}

func insultMessage(users model.Users, adj, noun, verb string, random rand.Rand) string {
	descriptor := "a"
	switch adj[0] {
	case 'a', 'e', 'i', 'o', 'u':
		descriptor += "n"
	}
	insult := ""
	switch random.Intn(8) + 1 {
	case 1:
		insult = fmt.Sprintf("%s, you %s like %s %s %s. - %s", users.To, verb, descriptor, adj, noun, users.From)
	case 2:
		insult = fmt.Sprintf("When god made %s, his primary source of inspiration was %s %s %s.  - %s", users.To, descriptor, adj, noun, users.From)
	case 3:
		insult = fmt.Sprintf("%s's fetishes involve %s %s %s.  - %s", users.To, descriptor, adj, noun, users.From)
	case 4:
		insult = fmt.Sprintf("I don't know what makes %s so stupid, but it's probably because they're %s %s %s. - %s", users.To, descriptor, adj, noun, users.From)
	case 5:
		insult = fmt.Sprintf("%s, just %s you %s 4head. - %s", users.To, verb, adj, users.From)
	case 6:
		insult = fmt.Sprintf("%s's brain is so odd, that if doctors cracked their head open, they'd think it was a %s %s - %s", users.To, adj, noun, users.From)
	case 7:
		insult = fmt.Sprintf("I truly want to believe the world is a good place, but I'm constantly reminded that %s %ss like %s exists in it.  - %s", adj, noun, users.To, users.From)
	case 8:
		insult = fmt.Sprintf("%s is as ugly as %s %s %s - %s", users.To, descriptor, adj, noun, users.From)
	}

	return insult
}
