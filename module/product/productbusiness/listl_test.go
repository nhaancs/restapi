package productbusiness

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"restapi/module/product/productmodel"
	"testing"
)

func (s *businessSuite) TestBusinessSuite_List() {
	req := productmodel.ListReq{
		Offset: 0,
		Limit:  1,
	}

	cases := []struct {
		name string
		req  productmodel.ListReq
		res  *productmodel.ListRes
		err  error
		mock func()
	}{
		{
			name: "success",
			req:  req,
			res:  &productmodel.ListRes{Products: []*productmodel.Product{{}}},
			err:  nil,
			mock: func() {
				s.productStore.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Return([]*productmodel.Product{{}}, nil).Times(1)
			},
		},
		{
			name: "can not list",
			req:  req,
			res:  nil,
			err:  errors.New("dummy error"),
			mock: func() {
				s.productStore.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("dummy error")).Times(1)
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			res, err := s.productBusiness.List(ctx, &c.req)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.res, res)
		})
	}
}
