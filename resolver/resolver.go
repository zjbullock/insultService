package resolver

import (
	"context"
	"randomInsultService/model"
	"randomInsultService/service"
)

type Resolver struct {
	Services struct {
		Insult service.Insult
	}
}

func (r *Resolver) GetInsult(ctx context.Context, args struct{ People *model.Users }) (*insultResolver, error) {
	insult := model.Insult{
		Message: "",
	}
	return &insultResolver{Insult: insult}, nil
}
