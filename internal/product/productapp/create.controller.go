package productapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) createController(c *gin.Context) {
	ctx := c.Request.Context()

	var req newProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		respond.Error(c, a.log, err)
		return
	}

	input, err := toBusNewProduct(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	output, err := a.productBus.Create(ctx, input)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "create: req[%+v]: %s", req, err))
		return
	}

	respond.Success(c, a.log, toAppProduct(output))
}
