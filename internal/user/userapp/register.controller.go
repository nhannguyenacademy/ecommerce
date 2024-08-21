package userapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

func (a *app) registerController(c *gin.Context) {
	ctx := c.Request.Context()

	var req registerUser
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	nu, err := toBusRegisterUser(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	usr, err := a.userBus.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			respond.Error(c, a.log, errs.New(errs.Aborted, userbus.ErrUniqueEmail))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "register: usr[%+v]: %s", usr, err))
		}
		return
	}
}
