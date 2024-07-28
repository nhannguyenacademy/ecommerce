package response

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"net/http"
)

type httpStatus interface {
	HTTPStatus() int
}

func Send(c *gin.Context, log *logger.Logger, data any, err error) error {
	// If the context has been canceled, it means the client is no longer waiting for a response.
	ctx := c.Request.Context()
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info(ctx, "client disconnected, do not send response")
		}
	}

	// handle errors
	if err != nil {
		if cErr := c.Error(err); cErr != nil {
			log.Warn(ctx, "send response: failed to add error to gin context", "error", cErr, "original error", err)
		}
		appErr := errs.NewError(err)
		c.AbortWithStatusJSON(appErr.HTTPStatus(), appErr)
		return nil
	}

	// handle data
	var statusCode = http.StatusOK
	switch v := data.(type) {
	case httpStatus:
		statusCode = v.HTTPStatus()
	case error:
		statusCode = http.StatusInternalServerError
		log.Warn(ctx, "send response: should not return error in data", "error", v)
	default:
		if data == nil {
			statusCode = http.StatusNoContent
		}
	}

	c.JSON(statusCode, data)

	return nil
}
