package productapp

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) deleteController(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid product id: %s", err))
		return
	}

	if err := a.productBus.Delete(ctx, id); err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "delete: id[%s]: %s", id, err))
		return
	}

	respond.Success(c, a.log, nil)
}
