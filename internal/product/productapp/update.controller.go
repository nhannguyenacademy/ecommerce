package productapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
)

func (a *app) updateController(c *gin.Context) {
	var req updateProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	u, err := a.update(c.Request.Context(), req, c.Param("product_id"))
	response.Send(c, a.log, u, err)
}

func (a *app) update(ctx context.Context, req updateProduct, id string) (product, error) {
	up, err := toBusUpdateProduct(req)
	if err != nil {
		return product{}, errs.New(errs.InvalidArgument, err)
	}

	prdID, err := uuid.Parse(id)
	if err != nil {
		return product{}, errs.Newf(errs.InvalidArgument, "invalid product id: %s", err)
	}

	prd, err := a.productBus.QueryByID(ctx, prdID)
	if err != nil {
		if errors.Is(err, productbus.ErrNotFound) {
			return product{}, errs.Newf(errs.NotFound, "querybyid: %s", err)
		}
		return product{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	updPrd, err := a.productBus.Update(ctx, prd, up)
	if err != nil {
		return product{}, errs.Newf(errs.Internal, "update: prdID[%s] up[%+v]: %s", prd.ID, up, err)
	}

	return toAppProduct(updPrd), nil
}
