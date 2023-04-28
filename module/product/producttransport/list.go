package producttransport

import (
	"github.com/gin-gonic/gin"
	"restapi/module/product/productmodel"
	"restapi/pkg/logging"
	"restapi/pkg/response"
)

func (t *transport) List() func(*gin.Context) {
	return func(c *gin.Context) {
		offset := c.GetInt("id")
		limit := c.GetInt("id")
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
