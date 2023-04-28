package productbusiness

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"restapi/module/product/productmodel"
	"testing"
)

func (s *businessSuite) TestBusinessSuite_Detail() {
	req := productmodel.DetailReq{
		Id: 1,
	}

	cases := []struct {
		name string
		req  productmodel.DetailReq
		res  *productmodel.DetailRes
		err  error
		mock func()
	}{
		{
			name: "success",
			req:  req,
			res:  &productmodel.DetailRes{Product: &productmodel.Product{}},
			err:  nil,
			mock: func() {
				s.productStore.On("FindById", mock.Anything, mock.Anything).
					Return(&productmodel.Product{}, nil).Times(1)
			},
		},
		{
			name: "can not get",
			req:  req,
			res:  nil,
			err:  errors.New("dummy error"),
			mock: func() {
				s.productStore.On("FindById", mock.Anything, mock.Anything).
					Return(nil, errors.New("dummy error")).Times(1)
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			res, err := s.productBusiness.Detail(ctx, &c.req)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.res, res)
		})
	}
}
