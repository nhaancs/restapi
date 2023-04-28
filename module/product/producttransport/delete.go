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

func (t *transport) Delete() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil || idInt <= 0 {
			response.Error(c, apperr.Wrap(err, appconst.CodeBadRequest, "invalid id", http.StatusBadRequest))
			return
		}

		logging.FromContext(c.Request.Context()).Info("start doing business logic")
		res, err := t.productBusiness.Delete(c.Request.Context(), &productmodel.DeleteReq{Id: idInt})
		if err != nil {
			logging.FromContext(c.Request.Context()).Errorf("error: %+v", err)
			response.Error(c, err)
			return
		}

		response.Success(c, res)
	}
}
