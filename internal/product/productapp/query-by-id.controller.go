package productapp

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
)

func (a *app) queryByIDController(c *gin.Context) {
	u, err := a.queryByID(c.Request.Context())
	response.Send(c, a.log, u, err)
}

func (a *app) queryByID(ctx context.Context) (user, error) {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return user{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppUser(usr), nil
}
