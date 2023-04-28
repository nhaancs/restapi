package productbusiness

import (
	"context"
	"restapi/module/product/productmodel"
)

type (
	ProductStore interface {
		Insert(ctx context.Context, product *productmodel.Product) error
		Select(ctx context.Context, offset, limit int) ([]*productmodel.Product, error)
		Update(ctx context.Context, product *productmodel.Product) error
		FindById(ctx context.Context, id int64) (*productmodel.Product, error)
		Delete(ctx context.Context, id int64) error
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
