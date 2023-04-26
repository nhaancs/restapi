package usertransport

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restapi/module/user/usermodel"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/response"
)

func (t *transport) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// parse request body
		var req usermodel.LoginReq
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, "bad request", http.StatusBadRequest))
			return
		}

		// validate request
		if err := req.Validate(); err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, err.Error(), http.StatusBadRequest))
			return
		}

		// do business logic
		res, err := t.userBusiness.Login(c.Request.Context(), &req)
		if err != nil {
			t.log.Errorf("login error: %+v", err)
			response.Error(c, err)
			return
		}

		response.Success(c, res)
	}
}
