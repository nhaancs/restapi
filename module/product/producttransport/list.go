package producttransport

import (
	"github.com/gin-gonic/gin"
	"restapi/module/product/productmodel"
	"restapi/pkg/logging"
	"restapi/pkg/response"
	"strconv"
)

func (t *transport) List() func(*gin.Context) {
	return func(c *gin.Context) {
		offsetStr := c.Query("offset")
		offset, _ := strconv.Atoi(offsetStr)

		limitStr := c.Query("limit")
		limit, _ := strconv.Atoi(limitStr)

		// validate request
		if limit <= 0 {
			limit = 20
		}

		logging.FromContext(c.Request.Context()).Info("start doing business logic", limit, offset)
		res, err := t.productBusiness.List(c.Request.Context(), &productmodel.ListReq{Offset: offset, Limit: limit})
		if err != nil {
			logging.FromContext(c.Request.Context()).Errorf("error: %+v", err)
			response.Error(c, err)
			return
		}

		response.Success(c, res)
	}
}
