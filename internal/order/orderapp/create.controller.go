package orderapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
)

func (a *app) createController(c *gin.Context) {
	var req newOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	u, err := a.create(c.Request.Context(), req)
	response.Send(c, a.log, u, err)
}

func (a *app) create(ctx context.Context, req newOrder) (order, error) {
	a, err := a.newWithTx(ctx)
	if err != nil {
		return order{}, errs.New(errs.Internal, err)
	}

	no, err := toBusNewOrder(req)
	if err != nil {
		return order{}, errs.New(errs.InvalidArgument, err)
	}

	ord, err := a.orderBus.Create(ctx, no)
	if err != nil {
		return order{}, errs.Newf(errs.Internal, "create: ord[%+v]: %s", ord, err)
	}

	return toAppOrder(ord), nil
}
