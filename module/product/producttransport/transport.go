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
		Update(ctx context.Context, req *productmodel.UpdateReq) (*productmodel.UpdateRes, error)
		List(ctx context.Context, req *productmodel.ListReq) (*productmodel.ListRes, error)
		Detail(ctx context.Context, req *productmodel.DetailReq) (*productmodel.DetailRes, error)
		Delete(ctx context.Context, req *productmodel.DeleteReq) (*productmodel.DeleteRes, error)
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
	users.GET(":id", t.Detail())
	users.GET("", t.List())
	users.POST("", middleware.Auth(t.tokenProvider), t.Create())
	users.PUT(":id", middleware.Auth(t.tokenProvider), t.Update())
	users.DELETE(":id", middleware.Auth(t.tokenProvider), t.Delete())
}
