package productbusiness

import (
	"context"
	"restapi/module/product/productmodel"
)

func (b *business) Delete(ctx context.Context, req *productmodel.DeleteReq) (*productmodel.DeleteRes, error) {
	err := b.productStore.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &productmodel.DeleteRes{}, nil
}
