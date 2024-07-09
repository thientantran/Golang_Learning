package userbiz

import (
	"Food-delivery/common"
	"Food-delivery/component/tokenprovider"
	usermodel "Food-delivery/module/user/model"
	"context"
)

type LoginStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

type loginBussiness struct {
	//appCtx        appctx.AppContext
	storeUser     LoginStore
	tokenProvider tokenprovider.Provider
	hasher        Hasher
	expiry        int
}

func NewLoginBussiness(storeUser LoginStore, tokenProvider tokenprovider.Provider, hasher Hasher, expiry int) *loginBussiness {
	return &loginBussiness{storeUser: storeUser, tokenProvider: tokenProvider, hasher: hasher, expiry: expiry}
}

// 1. Find user by email
// 2. Hash password from input and compare with password in database
// 3. Provider: Issue JWT token for client
// 3.1 Access token and refresh token
// 4. Return user info and token
func (bussiness *loginBussiness) Login(ctx context.Context, data *usermodel.UserLogin) (*tokenprovider.Token, error) {
	user, err := bussiness.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, usermodel.ErrEmailOrPasswordInvalid
	}
	passHashed := bussiness.hasher.Hash(data.Password + user.Salt)

	if user.Password != passHashed {
		return nil, usermodel.ErrEmailOrPasswordInvalid
	}

	payload := tokenprovider.TokenPayload{
		UserID: user.Id,
		Role:   user.Role,
	}

	accessToken, err := bussiness.tokenProvider.Generate(payload, bussiness.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	return accessToken, nil
}
