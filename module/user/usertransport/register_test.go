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
	"testing"
)

func (s *transportSuite) TestTransportSuite_Register() {
	registerReq := usermodel.RegisterReq{
		Username: "username",
		Password: "password",
	}
	cases := []struct {
		name    string
		args    func() io.Reader
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: func() io.Reader {
				data, _ := json.Marshal(registerReq)
				return bytes.NewReader(data)
			},
			wantErr: false,
			mock: func() {
				s.userBusiness.On("Register", mock.Anything, mock.Anything).Return(&usermodel.RegisterRes{}, nil).Times(1)
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
			req, _ := http.NewRequest(http.MethodPost, "/v1/users/register", c.args())
			router.ServeHTTP(w, req)
			if c.wantErr {
				assert.NotEqual(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
			assert.NotNil(t, w.Body.String())
		})
	}
}
