package mid

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"io"
	"net/url"
	"time"
)

// customResponseWriter wraps gin.ResponseWriter to capture the response body
type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logging(l *logger.Logger) gin.HandlerFunc {
	ginPath := func(c *gin.Context) string {
		if c.FullPath() != "" {
			return c.FullPath()
		}
		return c.Request.URL.Path
	}

	readBody := func(reader io.Reader) string {
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)

		s := buf.String()
		s, err := url.QueryUnescape(s)
		if err != nil {
			return fmt.Sprintf("readBody err: %s", err.Error())
		}
		return s
	}

	getReqBody := func(c *gin.Context) string {
		buf, _ := io.ReadAll(c.Request.Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		//We have to create a new Buffer, because rdr1 will be read.
		rdr2 := io.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = rdr2
		return readBody(rdr1)
	}

	return func(c *gin.Context) {
		ctx, start := c.Request.Context(), time.Now()

		l := l.With(map[string]any{
			"request_path":   ginPath(c),
			"remote_address": c.Request.RemoteAddr,
		})

		l.Info(
			ctx,
			"http server: received request",
			"method", c.Request.Method,
			"client_ip", c.ClientIP(),
			"request_body", getReqBody(c),
			"params", c.Request.URL.Query().Encode(),
		)

		// Replace the default ResponseWriter with the custom one
		customWriter := &customResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = customWriter
		c.Next()

		l = l.With(map[string]interface{}{
			"status":        c.Writer.Status(),
			"latency":       time.Since(start).String(),
			"response_body": customWriter.body.String(),
		})

		if len(c.Errors) > 0 {
			l = l.With(map[string]interface{}{
				"errors": c.Errors,
			})
		}

		l.Info(ctx, "http server: sending response")
	}
}
