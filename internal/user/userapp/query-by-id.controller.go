package userapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) queryByIDController(c *gin.Context) {
	ctx := c.Request.Context()

	usr, err := mid.GetUser(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "querybyid: %s", err))
		return
	}

	respond.Success(c, a.log, toAppUser(usr))
}
