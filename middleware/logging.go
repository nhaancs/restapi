package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/url"
	"restapi/pkg/tracing"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logging(log *zap.SugaredLogger) func(c *gin.Context) {
	ginPath := func(c *gin.Context) string {
		if c.FullPath() != "" {
			return c.FullPath()
		}
		return c.Request.URL.Path
	}

	return func(c *gin.Context) {
		ctx, start := c.Request.Context(), time.Now()
		l := log.With("trace_id", tracing.TraceID(ctx), "http.server.request_path", ginPath(c))
		body := getReqBody(c)

		l.With("http.server.method", c.Request.Method, "http.server.client_ip", c.ClientIP(), "payload", body).
			Infof("server: REQUEST")

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		l = l.With("http.server.status", c.Writer.Status(), "http.server.latency", time.Since(start).String(), "payload", blw.body.String())

		if len(c.Errors) > 0 {
			l = l.With("error", c.Errors)
		}

		if c.Writer.Status() >= 500 {
			l.Errorf("http.server: RESPONSE WITH ERROR")
		} else {
			l.Infof("http.server: RESPONSE")
		}
	}
}

func getReqBody(c *gin.Context) string {
	buf, _ := io.ReadAll(c.Request.Body)
	rdr1 := io.NopCloser(bytes.NewBuffer(buf))
	rdr2 := io.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
	c.Request.Body = rdr2
	return readBody(rdr1)
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	s, err := url.QueryUnescape(s)
	if err != nil {
		return fmt.Sprintf("readBody err: %s", err.Error())
	}
	return s
}
