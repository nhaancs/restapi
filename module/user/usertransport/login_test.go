package usertransport

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"restapi/module/user/usermodel"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"strings"
	"testing"
)

func (s *transportSuite) TestTransportSuite_Login() {
	loginReq := usermodel.LoginReq{
		Username: "username",
		Password: "password",
	}
	cases := []struct {
		name   string
		args   func() io.Reader
		status int
		mock   func()
	}{
		{
			name: "success",
			args: func() io.Reader {
				data, _ := json.Marshal(loginReq)
				return bytes.NewReader(data)
			},
			status: http.StatusOK,
			mock: func() {
				s.userBusiness.On("Login", mock.Anything, mock.Anything).Return(&usermodel.LoginRes{}, nil).Times(1)
			},
		},
		{
			name: "can not parse request body",
			args: func() io.Reader {
				return strings.NewReader("invalid")
			},
			status: http.StatusBadRequest,
			mock: func() {
			},
		},
		{
			name: "missing username",
			args: func() io.Reader {
				invalidReq := loginReq
				invalidReq.Username = ""
				data, _ := json.Marshal(invalidReq)
				return bytes.NewReader(data)
			},
			status: http.StatusBadRequest,
			mock: func() {
			},
		},
		{
			name: "missing password",
			args: func() io.Reader {
				invalidReq := loginReq
				invalidReq.Password = ""
				data, _ := json.Marshal(invalidReq)
				return bytes.NewReader(data)
			},
			status: http.StatusBadRequest,
			mock: func() {
			},
		},
		{
			name: "failed to do login business",
			args: func() io.Reader {
				data, _ := json.Marshal(loginReq)
				return bytes.NewReader(data)
			},
			status: http.StatusBadRequest,
			mock: func() {
				s.userBusiness.On("Login", mock.Anything, mock.Anything).
					Return(nil, apperr.Wrap(nil, appconst.CodeBadRequest, "", http.StatusBadRequest)).Times(1)
			},
		},
	}

	t := s.T()
	router := gin.New()
	s.userTransport.SetupRoutes(router.Group("v1"))

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/v1/users/login", c.args())
			router.ServeHTTP(w, req)

			assert.Equal(t, c.status, w.Code)
			assert.NotNil(t, w.Body.String())
		})
	}
}
