package productbusiness

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"restapi/module/product/productmodel"
	"testing"
)

func (s *businessSuite) TestBusinessSuite_Delete() {
	req := productmodel.DeleteReq{
		Id: 1,
	}

	cases := []struct {
		name string
		req  productmodel.DeleteReq
		res  *productmodel.DeleteRes
		err  error
		mock func()
	}{
		{
			name: "success",
			req:  req,
			res:  &productmodel.DeleteRes{},
			err:  nil,
			mock: func() {
				s.productStore.On("Delete", mock.Anything, mock.Anything).
					Return(nil).Times(1)
			},
		},
		{
			name: "can not delete",
			req:  req,
			res:  nil,
			err:  errors.New("dummy error"),
			mock: func() {
				s.productStore.On("Delete", mock.Anything, mock.Anything).
					Return(errors.New("dummy error")).Times(1)
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			res, err := s.productBusiness.Delete(ctx, &c.req)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.res, res)
		})
	}
}
