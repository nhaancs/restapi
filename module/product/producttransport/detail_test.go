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

func (s *transportSuite) TestTransportSuite_Detail() {
	cases := []struct {
		name    string
		status  int
		product *productmodel.Product
		mock    func()
	}{
		{
			name:    "success",
			status:  http.StatusOK,
			product: &productmodel.Product{},
			mock: func() {
				s.productBusiness.On("Detail", mock.Anything, mock.Anything, mock.Anything).
					Return(&productmodel.DetailRes{Product: &productmodel.Product{}}, nil).Times(1)
			},
		},
		{
			name:    "failed to do do business",
			status:  http.StatusBadRequest,
			product: nil,
			mock: func() {
				s.productBusiness.On("Detail", mock.Anything, mock.Anything, mock.Anything).
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
			req, _ := http.NewRequest(http.MethodGet, "/v1/products/1", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, c.status, w.Code)
			assert.NotNil(t, w.Body.String())
		})
	}

	t.Run("invalid product id", func(t *testing.T) {
		s.tokenProvider.On("Validate", mock.Anything).Return(nil, nil).Times(1)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/v1/products/abc", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotNil(t, w.Body.String())
	})
}
