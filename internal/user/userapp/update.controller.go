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
	var req updateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	u, err := a.update(c.Request.Context(), req)
	response.Send(c, a.log, u, err)
}

func (a *app) update(ctx context.Context, req updateUser) (user, error) {
	uu, err := toBusUpdateUser(req)
	if err != nil {
		return user{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return user{}, errs.Newf(errs.Internal, "user missing in context: %s", err)
	}

	updUsr, err := a.userBus.Update(ctx, usr, uu)
	if err != nil {
		return user{}, errs.Newf(errs.Internal, "update: userID[%s] uu[%+v]: %s", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}
