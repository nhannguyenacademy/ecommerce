package userapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

func (a *app) registerController(c *gin.Context) {
	var ru registerUser
	if err := c.ShouldBindJSON(&ru); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	u, err := a.register(c.Request.Context(), ru)
	response.Send(c, a.log, u, err)
}

func (a *app) register(ctx context.Context, ru registerUser) (user, error) {
	bru, err := toBusRegisterUser(ru)
	if err != nil {
		return user{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Register(ctx, bru)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return user{}, errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return user{}, errs.Newf(errs.Internal, "register: usr[%+v]: %s", usr, err)
	}

	return toAppUser(usr), nil
}
