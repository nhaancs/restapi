package productbusiness

import (
	"context"
	"restapi/module/product/productmodel"
)

func (b *business) Update(ctx context.Context, req *productmodel.UpdateReq) (*productmodel.UpdateRes, error) {
	err := b.productStore.Update(ctx, &productmodel.Product{Name: req.Name, Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &productmodel.UpdateRes{}, nil
}
