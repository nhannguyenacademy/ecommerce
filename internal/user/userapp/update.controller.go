package userapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
)

func (a *app) updateController(c *gin.Context) {
	ctx := c.Request.Context()

	var req updateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		respond.Error(c, a.log, err)
		return
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "user missing in context: %s", err))
		return
	}

	input, err := toBusUpdateUser(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	output, err := a.userBus.Update(ctx, usr, input)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "update: userID[%s] req[%+v]: %s", usr.ID, req, err))
		return
	}

	respond.Success(c, a.log, toAppUser(output))
}
