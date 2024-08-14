package productapp

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
)

func (a *app) queryByIDController(c *gin.Context) {
	u, err := a.queryByID(c.Request.Context(), c.Param("product_id"))
	response.Send(c, a.log, u, err)
}

func (a *app) queryByID(ctx context.Context, id string) (product, error) {
	prdID, err := uuid.Parse(id)
	if err != nil {
		return product{}, errs.Newf(errs.InvalidArgument, "invalid product id: %s", err)
	}
	usr, err := a.productBus.QueryByID(ctx, prdID)
	if err != nil {
		return product{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppProduct(usr), nil
}
