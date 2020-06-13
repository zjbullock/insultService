package resolver

import (
	"context"
	"github.com/juju/loggo"
	"github.com/pkg/errors"
	"insultService/model"
	"insultService/service"
	"insultService/service/helper"
)

// Resolver is a struct that contains different Services to be used when resolving a graphql query or mutation
type Resolver struct {
	Services struct {
		Insult service.Insult
		Sms    service.SMS
	}
	Log loggo.Logger
}

// GetInsult resolves an insult if any, and an error from the backend if present
func (r *Resolver) GetInsult(ctx context.Context, args struct{ People *model.Users }) (*insultResolver, error) {
	userInfo, err := r.Services.Insult.GetUserInfo(args.People.From)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting userInfo")
	}

	message, id, err := r.Services.Insult.GenerateInsult(*args.People, userInfo.Rank)
	if args.People.PhoneNumber != nil {
		secrets, err := helper.GetSecrets(r.Log)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting secrets")
		}
		r.Services.Sms.SendText(*args.People.PhoneNumber, *message, secrets)
	}
	promotion, err := r.Services.Insult.IncreaseUserExperience(args.People.From)

	insult := model.Insult{
		Message:   message,
		Promotion: promotion,
		Id:        id,
	}
	return &insultResolver{Insult: insult}, err
}

func (r *Resolver) GetUserInfo(ctx context.Context, args struct{ UserID string }) (*userInfoResolver, error) {
	userInfo, err := r.Services.Insult.GetUserInfo(args.UserID)
	r.Log.Infof("userInfo: %v", userInfo)
	return &userInfoResolver{UserInfo: *userInfo}, err
}

func (r *Resolver) GetInsultStats(ctx context.Context) (*insultStatResolver, error) {
	insultStat, err := r.Services.Insult.GetInsultsStats()
	if err != nil {
		return nil, err
	}
	return &insultStatResolver{InsultStat: *insultStat}, nil
}
