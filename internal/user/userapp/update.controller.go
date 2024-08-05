package userapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
)

func (a *app) updateController(c *gin.Context) {
	var uu UpdateUser
	if err := c.ShouldBindJSON(&uu); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	u, err := a.update(c.Request.Context(), uu)
	response.Send(c, a.log, u, err)
}

func (a *app) update(ctx context.Context, app UpdateUser) (User, error) {
	uu, err := toBusUpdateUser(app)
	if err != nil {
		return User{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return User{}, errs.Newf(errs.Internal, "user missing in context: %s", err)
	}

	updUsr, err := a.userBus.Update(ctx, usr, uu)
	if err != nil {
		return User{}, errs.Newf(errs.Internal, "update: userID[%s] uu[%+v]: %s", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}
