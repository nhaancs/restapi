package productbusiness

import (
	"context"
	"restapi/module/product/productmodel"
)

func (b *business) List(ctx context.Context, req *productmodel.ListReq) (*productmodel.ListRes, error) {
	products, err := b.productStore.Select(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	return &productmodel.ListRes{Products: products}, nil
}
