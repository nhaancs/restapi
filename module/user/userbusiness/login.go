package userbusiness

import (
	"context"
	"errors"
	"net/http"
	"restapi/module/user/usermodel"
	"restapi/module/user/userstore"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/logging"
	"restapi/pkg/tokenprovider"
)

func (b *business) Login(ctx context.Context, req *usermodel.LoginReq) (*usermodel.LoginRes, error) {
	user, err := b.userStore.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, userstore.ErrNotFound) {
			return nil, apperr.Wrap(err, appconst.CodeBadRequest, "invalid username or password", http.StatusBadRequest)
		}

		return nil, err
	}

	logging.FromContext(ctx).Infof("got user: %+v", user)

	hashedPass := b.hasher.Hash(req.Password + user.Salt)
	if user.HashedPassword != hashedPass {
		return nil, apperr.Wrap(nil, appconst.CodeBadRequest, "invalid username or password", http.StatusBadRequest)
	}

	token, err := b.tokenProvider.Generate(tokenprovider.TokenPayload{UserId: user.Id}, 365*24*3600)
	if err != nil {
		return nil, err
	}

	return &usermodel.LoginRes{
		Token: token.Token,
	}, nil
}
