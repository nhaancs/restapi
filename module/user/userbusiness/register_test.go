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
	"testing"
)

func (s *businessSuite) TestBusinessSuite_Register() {
	req := usermodel.RegisterReq{
		Username: "username",
		Password: "password",
	}

	cases := []struct {
		name string
		req  usermodel.RegisterReq
		res  *usermodel.RegisterRes
		err  error
		mock func()
	}{
		{
			name: "success",
			req:  req,
			res:  &usermodel.RegisterRes{},
			err:  nil,
			mock: func() {
				s.userStore.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, userstore.ErrNotFound).Times(1)
				s.hasher.On("Hash", mock.AnythingOfType("string")).
					Return("hashedpassword").Times(1)
				s.userStore.On("Insert", mock.Anything, mock.AnythingOfType("*usermodel.User")).
					Return(nil).Times(1)
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
			name: "duplicated username",
			req:  req,
			res:  nil,
			err:  apperr.Wrap(nil, appconst.CodeBadRequest, "username is already existed", http.StatusBadRequest),
			mock: func() {
				s.userStore.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(&usermodel.User{}, nil).Times(1)
			},
		},
		{
			name: "can not insert new user to database",
			req:  req,
			res:  nil,
			err:  errors.New("dummy error"),
			mock: func() {
				s.userStore.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, userstore.ErrNotFound).Times(1)
				s.hasher.On("Hash", mock.AnythingOfType("string")).
					Return("hashedpassword").Times(1)
				s.userStore.On("Insert", mock.Anything, mock.AnythingOfType("*usermodel.User")).
					Return(errors.New("dummy error")).Times(1)
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			res, err := s.userBusiness.Register(ctx, &c.req)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.res, res)
		})
	}
}
