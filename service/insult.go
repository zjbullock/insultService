package service

import (
	"fmt"
	"github.com/juju/loggo"
	"math/rand"
	"randomInsultService/model"
	"randomInsultService/repository"
	"time"
)

// Insult is an interface that contains methods relating to insults
type Insult interface {
	GenerateInsult(who model.Users) (message *string, id *string, err error)
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
	fmt.Printf("here are the words retrieved: %v\n", words)
	insultContents := model.InsultContent{
		Verb:      verb,
		Adjective: adj,
		Noun:      noun,
	}
	insult := insultMessage(who, adj, noun, verb)
	//Should insert generated insult into firebase collection
	id, err = i.fireStore.InsertEntry(insultContents)
	//Should produce an error if failed insert, but still return proper insult
	if err != nil {
		return &insult, nil, err
	}

	return &insult, id, nil
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
	insult := fmt.Sprintf("%s, you %s like %s %s %s. - %s", users.To, verb, descriptor, adj, noun, users.From)
	return insult
}
