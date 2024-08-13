package productapp

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
)

func (a *app) deleteController(c *gin.Context) {
	err := a.delete(c.Request.Context())
	response.Send(c, a.log, nil, err)
}

func (a *app) delete(ctx context.Context) error {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "userID missing in context: %s", err)
	}

	if err := a.userBus.Delete(ctx, usr); err != nil {
		return errs.Newf(errs.Internal, "delete: userID[%s]: %s", usr.ID, err)
	}

	return nil
}
