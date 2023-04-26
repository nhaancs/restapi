package userbusiness

import (
	"context"
	"go.uber.org/zap"
	"restapi/module/user/usermodel"
	"restapi/pkg/tokenprovider"
)

type (
	UserStore interface {
		FindByUsername(ctx context.Context, username string) (*usermodel.User, error)
		Insert(ctx context.Context, user *usermodel.User) error
	}
	TokenProvider interface {
		Generate(payload tokenprovider.TokenPayload, expiryInSeconds int64) (*tokenprovider.Token, error)
		Validate(token string) (*tokenprovider.TokenPayload, error)
	}
	Hasher interface {
		Hash(raw string) string
	}
)

type business struct {
	log           *zap.SugaredLogger
	userStore     UserStore
	tokenProvider TokenProvider
	hasher        Hasher
}

func New(
	log *zap.SugaredLogger,
	userStore UserStore,
	tokenProvider TokenProvider,
	hasher Hasher,
) *business {
	return &business{
		log:           log,
		userStore:     userStore,
		tokenProvider: tokenProvider,
		hasher:        hasher,
	}
}
