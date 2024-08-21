package productapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) updateController(c *gin.Context) {
	ctx := c.Request.Context()

	var req updateProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	id, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid id: %s", err))
		return
	}

	input, err := toBusUpdateProduct(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	output, err := a.productBus.Update(ctx, id, input)
	if err != nil {
		if errors.Is(err, productbus.ErrNotFound) {
			respond.Error(c, a.log, errs.Newf(errs.NotFound, "update: id[%s] req[%+v]: %s", id, req, err))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "update: id[%s] req[%+v]: %s", id, req, err))
		}
		return
	}

	respond.Success(c, a.log, toAppProduct(output))
}
