package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/metric"
	"strconv"
)

type Response struct {
	Code    string      `json:"error_code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Error(c *gin.Context, err error) {
	appErr := apperr.Convert(err)
	metric.Server().Inc(c.FullPath(), strconv.Itoa(appErr.Status()), appErr.Code())

	c.Header("Content-Type", "application/json")
	c.AbortWithStatusJSON(
		appErr.Status(),
		Response{
			Code:    appErr.Code(),
			Message: appErr.Message(),
		},
	)
}

func Success(c *gin.Context, data interface{}) {
	metric.Server().Inc(c.FullPath(), strconv.Itoa(http.StatusOK), appconst.CodeSuccess)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, Response{
		Code:    appconst.CodeSuccess,
		Message: "Success",
		Data:    data,
	})
}
