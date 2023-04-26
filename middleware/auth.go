package middleware

//
//import (
//	"context"
//	"errors"
//	usermodel "restapi/module/user/usermodel"
//	"restapi/pkg"
//	"restapi/pkg/apperr"
//	"restapi/pkg/appprovider"
//	"restapi/pkg/tokenprovider/jwtprovider"
//
//	"github.com/gin-gonic/gin"
//	"go.opencensus.io/trace"
//)
//
//type AuthenStore interface {
//	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
//}
//
//// 1. Get token from header
//// 2. Validate token and parse to payload
//// 3. From the token payload, we use user_id to find from DB
//func Auth(appCtx appprovider.Config, authStore AuthenStore) func(c *gin.Context) {
//	tokenProvider := jwtprovider.NewTokenJWTProvider(appCtx.SecretKey())
//	return func(c *gin.Context) {
//		token, err := pkg.ExtractTokenFromHeaderString(c.GetHeader("Authorization"))
//		if err != nil {
//			panic(err)
//		}
//
//		// database := appCtx.DB()
//		// userstore := userstorage.NewSQLStore(database)
//		payload, err := tokenProvider.Validate(token)
//		if err != nil {
//			panic(err)
//		}
//
//		ctx, span := trace.StartSpan(c.Request.Context(), "middleware.Auth.find-user")
//		user, err := authStore.FindUser(ctx, map[string]interface{}{"id": payload.UserId})
//		span.End()
//		if err != nil {
//			panic(err)
//		}
//
//		if user.DeletedAt != nil || !user.IsEnabled {
//			panic(apperr.ErrNoPermission(errors.New("user has been deleted or banned")))
//		}
//
//		user.Mask(user.Role == pkg.AdminRole)
//		c.Set(pkg.CurrentUser, user)
//		c.Next()
//	}
//}
