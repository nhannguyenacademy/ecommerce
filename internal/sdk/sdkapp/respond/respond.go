// Package respond contains functions to send responses to clients.
package respond

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"net/http"
)

type httpStatus interface {
	HTTPStatus() int
}

func Success(c *gin.Context, log *logger.Logger, data any) {
	ctx := c.Request.Context()

	var statusCode int
	switch v := data.(type) {
	case httpStatus:
		statusCode = v.HTTPStatus()
	case error:
		statusCode = http.StatusInternalServerError
		log.Warn(ctx, "respond success: should not return error in data", "error", v)
	default:
		statusCode = http.StatusOK
		if data == nil {
			statusCode = http.StatusNoContent
		}
	}

	c.JSON(statusCode, data)
}

func Error(c *gin.Context, log *logger.Logger, err error) {
	ctx := c.Request.Context()
	if err == nil {
		log.Warn(ctx, "respond error: error is nil")
	}

	if cErr := c.Error(err); cErr != nil {
		log.Warn(ctx, "respond error: failed to add error to gin context", "error", cErr, "original error", err)
	}

	var (
		appErr *errs.Error
		vErrs  validator.ValidationErrors
	)
	if errors.As(err, &vErrs) {
		appErr = errs.Newf(errs.InvalidArgument, "%s", vErrs)
	} else {
		appErr = errs.NewError(err)
	}

	c.AbortWithStatusJSON(appErr.HTTPStatus(), appErr)
}
