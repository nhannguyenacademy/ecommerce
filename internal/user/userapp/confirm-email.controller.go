package userapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

func (a *app) confirmEmailController(c *gin.Context) {
	ctx := c.Request.Context()

	confirmToken := c.Param("confirm_token")
	if confirmToken == "" {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, errors.New("missing confirm_token")))
	}

	err := a.userBus.ConfirmEmail(ctx, confirmToken)
	if err != nil {
		if errors.Is(err, userbus.ErrNotFound) {
			respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "confirmEmail: confirmToken[%s]: %s", confirmToken, err))
		}
		return
	}

	respond.Success(c, a.log, nil)
}
