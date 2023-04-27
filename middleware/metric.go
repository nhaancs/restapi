package middleware

import (
	"github.com/gin-gonic/gin"
	"restapi/pkg/metric"
)

func Metric() func(c *gin.Context) {
	return func(c *gin.Context) {
		method, path := c.Request.Method, c.FullPath()
		end := metric.HTTPServer().Start(method, path)
		c.Next()
		end(c.Writer.Status())
	}
}
