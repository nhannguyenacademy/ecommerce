package mid

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"runtime"
)

func Panic(l *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic: %+v - stack: %s", r, string(stack(3)))
				respond.Error(c, l, err)
				return
			}
		}()
		c.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	var buf = new(bytes.Buffer)
	for i := skip; ; i++ { // Skip the expected number of frames
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d\n", file, line)
	}
	return buf.Bytes()
}
