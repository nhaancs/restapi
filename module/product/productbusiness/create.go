package productbusiness

import (
	"context"
	"restapi/module/product/productmodel"
	"restapi/pkg/logging"
)

func (b *business) Create(ctx context.Context, req *productmodel.CreateReq) (*productmodel.CreateRes, error) {
	newProduct := productmodel.Product{Name: req.Name}
	logging.FromContext(ctx).Info("start insert new user to db")
	if err := b.productStore.Insert(ctx, &newProduct); err != nil {
		return nil, err
	}

	return &productmodel.CreateRes{}, nil
}
