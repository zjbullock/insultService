package service

import (
	"github.com/juju/loggo"
	"randomInsultService/model"
	"randomInsultService/repository"
)

type Insult interface {
	GenerateInsult(who model.Users) (string, error)
}

type insult struct {
	fireBase repository.FireBase
	log      loggo.Logger
}

func NewInsult(fire repository.FireBase, log loggo.Logger) Insult {
	return &insult{
		fireBase: fire,
		log:      log,
	}
}

func (i *insult) GenerateInsult(who model.Users) (string, error) {
	//Should generate an Insult
	//Should insert generated insult into firebase collection
	//Should produce an error if failed insert, but still return proper insult
	return "", nil
}
