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

func (s *transportSuite) TestTransportSuite_Register() {
	req := productmodel.CreateReq{
		Name: "product name",
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
				s.productBusiness.On("Create", mock.Anything, mock.Anything).Return(&productmodel.CreateRes{}, nil).Times(1)
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
			name: "missing product name",
			args: func() io.Reader {
				invalidReq := req
				invalidReq.Name = ""
				data, _ := json.Marshal(invalidReq)
				return bytes.NewReader(data)
			},
			status: http.StatusBadRequest,
			mock: func() {
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
				s.productBusiness.On("Create", mock.Anything, mock.Anything).
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
			req, _ := http.NewRequest(http.MethodPost, "/v1/products", c.args())
			router.ServeHTTP(w, req)

			assert.Equal(t, c.status, w.Code)
			assert.NotNil(t, w.Body.String())
		})
	}
}
