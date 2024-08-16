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
	var req registerUser
	if err := c.ShouldBindJSON(&req); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	err := a.register(c.Request.Context(), req)
	response.Send(c, a.log, nil, err)
}

func (a *app) register(ctx context.Context, req registerUser) error {
	nu, err := toBusRegisterUser(req)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return errs.Newf(errs.Internal, "register: usr[%+v]: %s", usr, err)
	}

	// todo: send email confirmation

	return nil
}
