package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
)

type Response struct {
	Code    string      `json:"error_code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Error(c *gin.Context, err error) {
	c.Header("Content-Type", "application/json")
	appErr := apperr.Convert(err)

	c.AbortWithStatusJSON(
		appErr.Status(),
		Response{
			Code:    appErr.Code(),
			Message: appErr.Message(),
		},
	)
}

func Success(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, Response{
		Code:    appconst.CodeSuccess,
		Message: "Success",
		Data:    data,
	})
}
