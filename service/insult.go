package service

import (
	"fmt"
	"github.com/juju/loggo"
	"insultService/model"
	"insultService/repository"
	"math/rand"
	"sort"
	"time"
)

// Insult is an interface that contains methods relating to insults
type Insult interface {
	GenerateInsult(who model.Users) (message *string, id *string, err error)
	GetInsultsStats() (insultStat *model.InsultStat, err error)
	GetUserInfo(userID string) (userInfo *model.UserInfo, err error)
	IncreaseUserExperience(userID string) (message *string, err error)
}

type insult struct {
	fireStore repository.FireStore
	log       loggo.Logger
}

// NewInsult creates a new insult service
func NewInsult(fire repository.FireStore, log loggo.Logger) Insult {
	return &insult{
		fireStore: fire,
		log:       log,
	}
}

// GenerateInsult returns a string with a generated insult and an error bubbled up from firestore if any
func (i *insult) GenerateInsult(who model.Users) (message *string, id *string, err error) {
	//Should generate an Insult
	words, err := i.fireStore.ReadAllWords()
	if err != nil {
		return nil, nil, err
	}
	adj, noun, verb := randomWordChooser(words)
	insultContents := model.InsultContent{
		Verb:      verb,
		Adjective: adj,
		Noun:      noun,
	}
	insult := insultMessage(who, adj, noun, verb)
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
		userInfo := model.UserInfo{
			Name:       userID,
			Rank:       determineTitle(1, *titles),
			Experience: 1,
		}
		i.log.Infof("userInfo: %v", userInfo)
		err = i.fireStore.UpdateUserInfo(&userInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to create user with error: %v", err)
		}
		message := fmt.Sprintf("Welcome!  Successfully Created a new user of the Insult Bot!  Your Current Rank is: %s", userInfo.Rank)
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
	userInfoList, err := i.fireStore.ReadUserInfo([]model.QueryArg{{Path: "username", Op: "==", Value: userID}})
	if err != nil {
		i.log.Errorf("error retrieving user with this username")
		return nil, err
	}
	if len(userInfoList) == 0 {
		return nil, fmt.Errorf("no such username exists")
	}
	return userInfoList[0], nil
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

func randomWordChooser(words *model.Words) (adjective, noun, verb string) {
	rand.Seed(time.Now().UTC().UnixNano())
	adjective = words.Adjective[rand.Intn(len(words.Adjective))]
	noun = words.Noun[rand.Intn(len(words.Noun))]
	verb = words.Verb[rand.Intn(len(words.Verb))]
	return adjective, noun, verb
}

func insultMessage(users model.Users, adj, noun, verb string) string {
	descriptor := "a"
	switch adj[0] {
	case 'a', 'e', 'i', 'o', 'u':
		descriptor += "n"
	}
	rand.Seed(time.Now().UTC().UnixNano())
	insult := ""
	switch rand.Intn(5) + 1 {
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
	}

	return insult
}
