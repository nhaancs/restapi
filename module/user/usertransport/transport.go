package usertransport

import (
	"context"
	"github.com/gin-gonic/gin"
	"restapi/module/user/usermodel"
)

type UserBusiness interface {
	Register(ctx context.Context, req *usermodel.RegisterReq) (*usermodel.RegisterRes, error)
	Login(ctx context.Context, req *usermodel.LoginReq) (*usermodel.LoginRes, error)
}

type transport struct {
	userBusiness UserBusiness
}

func New(
	userBusiness UserBusiness,
) *transport {
	return &transport{
		userBusiness: userBusiness,
	}
}

func (t *transport) SetupRoutes(r *gin.RouterGroup) {
	users := r.Group("users")
	users.POST("register", t.Register())
	users.POST("login", t.Login())
}
