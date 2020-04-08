package resolver

import "insultService/model"

type adjectiveCountResolver struct {
	AdjectiveCount model.AdjectiveCount
}

func (a *adjectiveCountResolver) Word() string {
	return a.AdjectiveCount.Adjective
}

func (a *adjectiveCountResolver) Count() int32 {
	return int32(a.AdjectiveCount.Count)
}
