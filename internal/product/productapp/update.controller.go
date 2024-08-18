package productapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) updateController(c *gin.Context) {
	ctx := c.Request.Context()

	var req updateProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		respond.Error(c, a.log, err)
		return
	}

	up, err := toBusUpdateProduct(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	prdID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid product id: %s", err))
		return
	}

	prd, err := a.productBus.QueryByID(ctx, prdID)
	if err != nil {
		var appErr *errs.Error
		if errors.Is(err, productbus.ErrNotFound) {
			appErr = errs.Newf(errs.NotFound, "querybyid: %s", err)
		} else {
			appErr = errs.Newf(errs.Internal, "querybyid: %s", err)
		}
		respond.Error(c, a.log, appErr)
		return
	}

	updPrd, err := a.productBus.Update(ctx, prd, up)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "update: prdID[%s] up[%+v]: %s", prd.ID, up, err))
		return
	}

	respond.Success(c, a.log, toAppProduct(updPrd))
}
