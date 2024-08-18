package productapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) deleteController(c *gin.Context) {
	ctx := c.Request.Context()

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

	if err := a.productBus.Delete(ctx, prd); err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "delete: prdID[%s]: %s", prd.ID, err))
		return
	}

	respond.Success(c, a.log, nil)
}
