package userbusiness

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"restapi/module/user/usermodel"
	"restapi/module/user/userstore"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/tokenprovider"
	"testing"
)

func (s *businessSuite) TestBusinessSuite_Login() {
	req := usermodel.LoginReq{
		Username: "username",
		Password: "password",
	}

	cases := []struct {
		name string
		req  usermodel.LoginReq
		res  *usermodel.LoginRes
		err  error
		mock func()
	}{
		{
			name: "success",
			req:  req,
			res:  &usermodel.LoginRes{Token: "token"},
			err:  nil,
			mock: func() {
				s.userStore.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(&usermodel.User{HashedPassword: "hashedpassword"}, nil).Times(1)
				s.hasher.On("Hash", mock.AnythingOfType("string")).
					Return("hashedpassword").Times(1)
				s.tokenProvider.On("Generate", mock.AnythingOfType("tokenprovider.TokenPayload"), mock.AnythingOfType("int64")).
					Return(&tokenprovider.Token{Token: "token"}, nil).Times(1)
			},
		},
		{
			name: "username not found",
			req:  req,
			res:  nil,
			err:  apperr.Wrap(userstore.ErrNotFound, appconst.CodeBadRequest, "invalid username or password", http.StatusBadRequest),
			mock: func() {
				s.userStore.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, userstore.ErrNotFound).Times(1)
			},
		},
		{
			name: "database internal error",
			req:  req,
			res:  nil,
			err:  errors.New("dummy error"),
			mock: func() {
				s.userStore.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("dummy error")).Times(1)
			},
		},
		{
			name: "wrong password",
			req:  req,
			res:  nil,
			err:  apperr.Wrap(nil, appconst.CodeBadRequest, "invalid username or password", http.StatusBadRequest),
			mock: func() {
				s.userStore.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(&usermodel.User{HashedPassword: "hashedpassword"}, nil).Times(1)
				s.hasher.On("Hash", mock.AnythingOfType("string")).
					Return("hashedpassword1").Times(1)
			},
		},
		{
			name: "can not gen jwt token",
			req:  req,
			res:  nil,
			err:  errors.New("dummy error"),
			mock: func() {
				s.userStore.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(&usermodel.User{HashedPassword: "hashedpassword"}, nil).Times(1)
				s.hasher.On("Hash", mock.AnythingOfType("string")).
					Return("hashedpassword").Times(1)
				s.tokenProvider.On("Generate", mock.AnythingOfType("tokenprovider.TokenPayload"), mock.AnythingOfType("int64")).
					Return(nil, errors.New("dummy error")).Times(1)
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			res, err := s.userBusiness.Login(ctx, &c.req)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.res, res)
		})
	}
}
