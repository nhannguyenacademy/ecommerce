package orderapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
)

func (a *app) deleteController(c *gin.Context) {
	err := a.delete(c.Request.Context(), c.Param("product_id"))
	response.Send(c, a.log, nil, err)
}

func (a *app) delete(ctx context.Context, id string) error {
	prdID, err := uuid.Parse(id)
	if err != nil {
		return errs.Newf(errs.InvalidArgument, "invalid product id: %s", err)
	}

	prd, err := a.productBus.QueryByID(ctx, prdID)
	if err != nil {
		if errors.Is(err, productbus.ErrNotFound) {
			return errs.Newf(errs.NotFound, "querybyid: %s", err)
		}
		return errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	if err := a.productBus.Delete(ctx, prd); err != nil {
		return errs.Newf(errs.Internal, "delete: prdID[%s]: %s", prd.ID, err)
	}

	return nil
}
