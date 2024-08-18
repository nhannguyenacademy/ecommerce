package orderapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) createController(c *gin.Context) {
	ctx := c.Request.Context()

	var req newOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		respond.Error(c, a.log, err)
		return
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.Internal, err))
		return
	}

	no, err := toBusNewOrder(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	ord, err := a.orderBus.Create(ctx, no)
	if err != nil {
		var appErr *errs.Error
		if errors.Is(err, orderbus.ErrMissingProducts) {
			appErr = errs.New(errs.InvalidArgument, orderbus.ErrMissingProducts)
		} else {
			appErr = errs.Newf(errs.Internal, "create: ord[%+v]: %s", ord, err)
		}
		respond.Error(c, a.log, appErr)
		return
	}

	respond.Success(c, a.log, toAppOrder(ord))
}
