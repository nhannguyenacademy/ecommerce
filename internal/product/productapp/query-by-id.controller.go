package productapp

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) queryByIDController(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid product id: %s", err))
		return
	}

	output, err := a.productBus.QueryByID(ctx, id)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "querybyid: %s", err))
		return
	}

	respond.Success(c, a.log, toAppProduct(output))
}
