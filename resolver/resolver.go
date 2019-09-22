package resolver

import (
	"context"
	"randomInsultService/model"
	"randomInsultService/service"
)

// Resolver is a struct that contains different Services to be used when resolving a graphql query or mutation
type Resolver struct {
	Services struct {
		Insult service.Insult
	}
}

// GetInsult resolves an insult if any, and an error from the backend if present
func (r *Resolver) GetInsult(ctx context.Context, args struct{ People *model.Users }) (*insultResolver, error) {
	message, id, err := r.Services.Insult.GenerateInsult(*args.People)
	if err != nil {
		return &insultResolver{
			Insult: model.Insult{
				Message: message,
				Id:      id,
			},
		}, err
	}
	insult := model.Insult{
		Message: "",
	}
	return &insultResolver{Insult: insult}, nil
}
