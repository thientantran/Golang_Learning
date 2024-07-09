package tokenprovider

import (
	"Food-delivery/common"
	"errors"
	"time"
)

type Token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

type TokenPayload struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}

type Provider interface {
	Generate(payload TokenPayload, expiry int) (*Token, error)
	Validate(token string) (*TokenPayload, error)
}

var (
	ErrNotFound = common.NewCustomError(
		errors.New("token not found"),
		"token not found",
		"ErrNotFound",
	)

	ErrEncodingToken = common.NewCustomError(
		errors.New("error encoding token"),
		"error encoding token",
		"ErrEncodingToken",
	)

	ErrInvalidToken = common.NewCustomError(
		errors.New("invalid token provided"),
		"invalid token provided",
		"ErrInvalidToken",
	)
)
