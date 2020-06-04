package resolver

import (
	"context"
	"insultService/model"
)

type insultStatResolver struct {
	InsultStat model.InsultStat
}

func (i *insultStatResolver) Verbs(_ context.Context) *[]*verbCountResolver {
	var verbs []*verbCountResolver
	for _, v := range i.InsultStat.Verbs {
		verbs = append(verbs, &verbCountResolver{
			VerbCount: model.VerbCount{
				Verb:  v.Verb,
				Count: v.Count,
			},
		})
	}
	return &verbs
}

func (i *insultStatResolver) Adjectives(_ context.Context) *[]*adjectiveCountResolver {
	var adjectives []*adjectiveCountResolver
	for _, a := range i.InsultStat.Adjectives {
		adjectives = append(adjectives, &adjectiveCountResolver{
			AdjectiveCount: model.AdjectiveCount{
				Adjective: a.Adjective,
				Count:     a.Count,
			},
		})
	}
	return &adjectives
}

func (i *insultStatResolver) Nouns(_ context.Context) *[]*nounCountResolver {
	var nouns []*nounCountResolver
	for _, v := range i.InsultStat.Nouns {
		nouns = append(nouns, &nounCountResolver{
			NounCount: model.NounCount{
				Noun:  v.Noun,
				Count: v.Count,
			},
		})
	}
	return &nouns
}
