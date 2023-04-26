package usertransport

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restapi/module/user/usermodel"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/response"
)

func (t *transport) Register() func(*gin.Context) {
	return func(c *gin.Context) {
		// parse request body
		var req usermodel.RegisterReq
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, "bad request", http.StatusBadRequest))
			return
		}

		// validate request
		if err := req.Validate(); err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, err.Error(), http.StatusBadRequest))
			return
		}

		res, err := t.userBusiness.Register(c.Request.Context(), &req)
		if err != nil {
			t.log.Errorf("register error: %+v", err)
			response.Error(c, err)
			return
		}

		response.Success(c, res)
	}
}
