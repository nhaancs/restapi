package producttransport

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"restapi/module/product/productmodel"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"testing"
)

func (s *transportSuite) TestTransportSuite_Delete() {
	cases := []struct {
		name   string
		status int
		mock   func()
	}{
		{
			name:   "success",
			status: http.StatusOK,
			mock: func() {
				s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
				s.productBusiness.On("Delete", mock.Anything, mock.Anything).Return(&productmodel.DeleteRes{}, nil).Times(1)
			},
		},
		{
			name:   "failed to do create business",
			status: http.StatusBadRequest,
			mock: func() {
				s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
				s.productBusiness.On("Delete", mock.Anything, mock.Anything).
					Return(nil, apperr.Wrap(nil, appconst.CodeBadRequest, "", http.StatusBadRequest)).Times(1)
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
			req, _ := http.NewRequest(http.MethodDelete, "/v1/products/1", nil)
			req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7InVzZXJfaWQiOjB9LCJleHAiOjE3MTQxNzMzNjYsImlhdCI6MTY4MjYzNzM2Nn0.H_RdlCzmiv17k3e9IYeciXoZg4N_pz59ECqVaiK5ORY")
			router.ServeHTTP(w, req)

			assert.Equal(t, c.status, w.Code)
			assert.NotNil(t, w.Body.String())
		})
	}

	t.Run("invalid product id", func(t *testing.T) {
		s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/v1/products/abc", nil)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7InVzZXJfaWQiOjB9LCJleHAiOjE3MTQxNzMzNjYsImlhdCI6MTY4MjYzNzM2Nn0.H_RdlCzmiv17k3e9IYeciXoZg4N_pz59ECqVaiK5ORY")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotNil(t, w.Body.String())
	})
}
