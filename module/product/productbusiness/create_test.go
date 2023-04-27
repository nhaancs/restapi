package productbusiness

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"restapi/module/product/productmodel"
	"testing"
)

func (s *businessSuite) TestBusinessSuite_Create() {
	req := productmodel.CreateReq{
		Name: "product name",
	}

	cases := []struct {
		name string
		req  productmodel.CreateReq
		res  *productmodel.CreateRes
		err  error
		mock func()
	}{
		{
			name: "success",
			req:  req,
			res:  &productmodel.CreateRes{},
			err:  nil,
			mock: func() {
				s.productStore.On("Insert", mock.Anything, mock.AnythingOfType("*productmodel.Product")).
					Return(nil).Times(1)
			},
		},
		{
			name: "can not insert new product to database",
			req:  req,
			res:  nil,
			err:  errors.New("dummy error"),
			mock: func() {
				s.productStore.On("Insert", mock.Anything, mock.AnythingOfType("*productmodel.Product")).
					Return(errors.New("dummy error")).Times(1)
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			res, err := s.productBusiness.Create(ctx, &c.req)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.res, res)
		})
	}
}
