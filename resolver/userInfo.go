package resolver

import (
	"context"
	"insultService/model"
)

type userInfoResolver struct {
	UserInfo model.UserInfo
}

func (u *userInfoResolver) Name(_ context.Context) string {
	return u.UserInfo.Name
}

func (u *userInfoResolver) Experience(_ context.Context) int32 {
	return int32(u.UserInfo.Experience)
}

func (u *userInfoResolver) Rank(_ context.Context) string {
	return u.UserInfo.Rank
}
