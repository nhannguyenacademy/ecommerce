package mid

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"runtime"
)

func Panic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.Error(errs.NewError(errors.New(fmt.Sprintf("panic: %+v - stack: %s", r, string(stack(3))))))
				appErr := errs.New(errs.Internal, errors.New("internal server error"))
				c.AbortWithStatusJSON(appErr.HTTPStatus(), appErr)
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
