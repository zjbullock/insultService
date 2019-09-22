package service

import (
	"github.com/juju/loggo"
	"randomInsultService/model"
	"randomInsultService/repository"
)

// Insult is an interface that contains methods relating to insults
type Insult interface {
	GenerateInsult(who model.Users) (message string, id *string, err error)
}

type insult struct {
	fireBase repository.FireStore
	log      loggo.Logger
}

// NewInsult creates a new insult service
func NewInsult(fire repository.FireStore, log loggo.Logger) Insult {
	return &insult{
		fireBase: fire,
		log:      log,
	}
}

// GenerateInsult returns a string with a generated insult and an error bubbled up from firebase if any
func (i *insult) GenerateInsult(who model.Users) (message string, id *string, err error) {
	//Should generate an Insult
	insultContents := model.InsultContent{
		Verb:      "taste",
		Adjective: "ugly",
		Noun:      "socks",
	}
	//Should insert generated insult into firebase collection
	id, err = i.fireBase.InsertEntry(insultContents)
	//Should produce an error if failed insert, but still return proper insult
	if err != nil {
		return "", nil, err
	}
	return "", id, nil
}
