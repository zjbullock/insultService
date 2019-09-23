package resolver

import (
	"context"
	"randomInsultService/model"
)

type insultResolver struct {
	Insult model.Insult
}

func (i *insultResolver) Message(_ context.Context) *string {
	return i.Insult.Message
}

func (i *insultResolver) FireStoreId(_ context.Context) *string {
	return i.Insult.Id
}
