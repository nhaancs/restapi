package producttransport

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restapi/module/product/productmodel"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/logging"
	"restapi/pkg/response"
)

func (t *transport) Create() func(*gin.Context) {
	return func(c *gin.Context) {
		// parse request body
		var req productmodel.CreateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, "bad request", http.StatusBadRequest))
			return
		}

		// validate request
		if err := req.Validate(); err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, err.Error(), http.StatusBadRequest))
			return
		}

		logging.FromContext(c.Request.Context()).Info("start doing business logic")
		res, err := t.productBusiness.Create(c.Request.Context(), &req)
		if err != nil {
			logging.FromContext(c.Request.Context()).Errorf("register error: %+v", err)
			response.Error(c, err)
			return
		}

		response.Success(c, res)
	}
}
