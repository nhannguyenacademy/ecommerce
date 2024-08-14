package productapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

func (a *app) createController(c *gin.Context) {
	var np newProduct
	if err := c.ShouldBindJSON(&np); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	u, err := a.create(c.Request.Context(), np)
	response.Send(c, a.log, u, err)
}

func (a *app) create(ctx context.Context, app newProduct) (product, error) {
	nc, err := toBusNewProduct(app)
	if err != nil {
		return product{}, errs.New(errs.InvalidArgument, err)
	}

	prd, err := a.productBus.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return product{}, errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return product{}, errs.Newf(errs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}
