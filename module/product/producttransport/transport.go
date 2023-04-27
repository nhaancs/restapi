package producttransport

import (
	"context"
	"github.com/gin-gonic/gin"
	"restapi/middleware"
	"restapi/module/product/productmodel"
	"restapi/pkg/tokenprovider"
)

type (
	ProductBusiness interface {
		Create(ctx context.Context, req *productmodel.CreateReq) (*productmodel.CreateRes, error)
	}
	TokenProvider interface {
		Validate(token string) (*tokenprovider.TokenPayload, error)
	}
)

type transport struct {
	productBusiness ProductBusiness
	tokenProvider   TokenProvider
}

func New(
	productBusiness ProductBusiness,
	tokenProvider TokenProvider,
) *transport {
	return &transport{
		productBusiness: productBusiness,
		tokenProvider:   tokenProvider,
	}
}

func (t *transport) SetupRoutes(r *gin.RouterGroup) {
	users := r.Group("products")
	users.POST("", middleware.Auth(t.tokenProvider), t.Create())
}
