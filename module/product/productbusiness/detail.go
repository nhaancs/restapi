package productbusiness

import (
	"context"
	"restapi/module/product/productmodel"
)

func (b *business) Detail(ctx context.Context, req *productmodel.DetailReq) (*productmodel.DetailRes, error) {
	product, err := b.productStore.FindById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &productmodel.DetailRes{Product: product}, nil
}
