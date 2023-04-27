package producttransport

import (
	"context"
	"github.com/gin-gonic/gin"
	"restapi/module/product/productmodel"
)

type ProductBusiness interface {
	Create(ctx context.Context, req *productmodel.CreateReq) (*productmodel.CreateRes, error)
}

type transport struct {
	productBusiness ProductBusiness
}

func New(
	productBusiness ProductBusiness,
) *transport {
	return &transport{
		productBusiness: productBusiness,
	}
}

func (t *transport) SetupRoutes(r *gin.RouterGroup) {
	users := r.Group("products")
	users.POST("", t.Create())
}
