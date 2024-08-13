package productapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

func (a *app) createController(c *gin.Context) {
	var nu newUser
	if err := c.ShouldBindJSON(&nu); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	u, err := a.create(c.Request.Context(), nu)
	response.Send(c, a.log, u, err)
}

func (a *app) create(ctx context.Context, app newUser) (user, error) {
	nc, err := toBusNewUser(app)
	if err != nil {
		return user{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return user{}, errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return user{}, errs.Newf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	return toAppUser(usr), nil
}
