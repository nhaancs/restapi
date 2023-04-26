package userbusiness

import (
	"context"
	"errors"
	"net/http"
	"restapi/module/user/usermodel"
	"restapi/module/user/userstore"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/salt"
)

func (b *business) Register(ctx context.Context, req *usermodel.RegisterReq) (*usermodel.RegisterRes, error) {
	user, err := b.userStore.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, userstore.ErrNotFound) {
		return nil, err
	}
	if user != nil {
		return nil, apperr.Wrap(err, appconst.CodeBadRequest, "username is already existed", http.StatusBadRequest)
	}

	salt := salt.GenSalt(50)
	hashedPass := b.hasher.Hash(req.Password + salt)
	newUser := usermodel.User{
		Username:       req.Username,
		HashedPassword: hashedPass,
		Salt:           salt,
	}

	if err := b.userStore.Insert(ctx, &newUser); err != nil {
		return nil, err
	}

	return &usermodel.RegisterRes{}, nil
}
