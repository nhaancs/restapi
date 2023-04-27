package productbusiness

import (
	"context"
	"restapi/module/product/productmodel"
)

type (
	ProductStore interface {
		Insert(ctx context.Context, product *productmodel.Product) error
	}
)

type business struct {
	productStore ProductStore
}

func New(
	productStore ProductStore,
) *business {
	return &business{
		productStore: productStore,
	}
}
