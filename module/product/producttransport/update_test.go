package producttransport

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"restapi/module/product/productmodel"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"strings"
	"testing"
)

func (s *transportSuite) TestTransportSuite_Update() {
	req := productmodel.UpdateReq{
		Name: "updated name",
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
				data, _ := json.Marshal(req)
				return bytes.NewReader(data)
			},
			status: http.StatusOK,
			mock: func() {
				s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
				s.productBusiness.On("Update", mock.Anything, mock.Anything).Return(&productmodel.UpdateRes{}, nil).Times(1)
			},
		},
		{
			name: "failed to do create business",
			args: func() io.Reader {
				data, _ := json.Marshal(req)
				return bytes.NewReader(data)
			},
			status: http.StatusBadRequest,
			mock: func() {
				s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
				s.productBusiness.On("Update", mock.Anything, mock.Anything).
					Return(nil, apperr.Wrap(nil, appconst.CodeBadRequest, "", http.StatusBadRequest)).Times(1)
			},
		},
		{
			name: "missing product name",
			args: func() io.Reader {
				invalidReq := req
				invalidReq.Name = ""
				data, _ := json.Marshal(invalidReq)
				return bytes.NewReader(data)
			},
			status: http.StatusBadRequest,
			mock: func() {
				s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
			},
		},
		{
			name: "can not parse request body",
			args: func() io.Reader {
				return strings.NewReader("invalid")
			},
			status: http.StatusBadRequest,
			mock: func() {
				s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
			},
		},
	}

	t := s.T()
	router := gin.New()
	s.productTransport.SetupRoutes(router.Group("v1"))

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/v1/products/1", c.args())
			req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7InVzZXJfaWQiOjB9LCJleHAiOjE3MTQxNzMzNjYsImlhdCI6MTY4MjYzNzM2Nn0.H_RdlCzmiv17k3e9IYeciXoZg4N_pz59ECqVaiK5ORY")
			router.ServeHTTP(w, req)

			assert.Equal(t, c.status, w.Code)
			assert.NotNil(t, w.Body.String())
		})
	}

	t.Run("invalid product id", func(t *testing.T) {
		s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
		w := httptest.NewRecorder()
		data, _ := json.Marshal(req)
		req, _ := http.NewRequest(http.MethodPut, "/v1/products/abc", bytes.NewReader(data))
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7InVzZXJfaWQiOjB9LCJleHAiOjE3MTQxNzMzNjYsImlhdCI6MTY4MjYzNzM2Nn0.H_RdlCzmiv17k3e9IYeciXoZg4N_pz59ECqVaiK5ORY")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotNil(t, w.Body.String())
	})
}
