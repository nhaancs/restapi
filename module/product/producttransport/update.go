package producttransport

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restapi/module/product/productmodel"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/logging"
	"restapi/pkg/response"
	"strconv"
)

func (t *transport) Update() func(*gin.Context) {
	return func(c *gin.Context) {
		// parse request body
		var req productmodel.UpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, "bad request", http.StatusBadRequest))
			return
		}

		// validate request
		if err := req.Validate(); err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, err.Error(), http.StatusBadRequest))
			return
		}

		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil || idInt <= 0 {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, "invalid id", http.StatusBadRequest))
			return
		}
		req.Id = idInt

		logging.FromContext(c.Request.Context()).Infof("start doing business logic: %+v", req)
		res, err := t.productBusiness.Update(c.Request.Context(), &req)
		if err != nil {
			logging.FromContext(c.Request.Context()).Errorf("error: %+v", err)
			response.Error(c, err)
			return
		}

		response.Success(c, res)
	}
}
