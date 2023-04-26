package usertransport

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"restapi/module/user/usermodel"
)

type UserBusiness interface {
	Register(ctx context.Context, req *usermodel.RegisterReq) (*usermodel.RegisterRes, error)
	Login(ctx context.Context, req *usermodel.LoginReq) (*usermodel.LoginRes, error)
}

type transport struct {
	log          *zap.SugaredLogger
	userBusiness UserBusiness
}

func New(
	log *zap.SugaredLogger,
	userBusiness UserBusiness,
) *transport {
	return &transport{
		log:          log,
		userBusiness: userBusiness,
	}
}

func (t *transport) SetupRoutes(r *gin.RouterGroup) {
	users := r.Group("users")
	users.POST("register", t.Register())
	users.POST("login", t.Login())
}
