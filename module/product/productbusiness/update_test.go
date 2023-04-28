package productbusiness

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"restapi/module/product/productmodel"
	"testing"
)

func (s *businessSuite) TestBusinessSuite_Update() {
	req := productmodel.UpdateReq{
		Id:   1,
		Name: "updated name",
	}

	cases := []struct {
		name string
		req  productmodel.UpdateReq
		res  *productmodel.UpdateRes
		err  error
		mock func()
	}{
		{
			name: "success",
			req:  req,
			res:  &productmodel.UpdateRes{},
			err:  nil,
			mock: func() {
				s.productStore.On("Update", mock.Anything, mock.Anything).
					Return(nil).Times(1)
			},
		},
		{
			name: "can not update",
			req:  req,
			res:  nil,
			err:  errors.New("dummy error"),
			mock: func() {
				s.productStore.On("Update", mock.Anything, mock.Anything).
					Return(errors.New("dummy error")).Times(1)
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			res, err := s.productBusiness.Update(ctx, &c.req)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.res, res)
		})
	}
}
