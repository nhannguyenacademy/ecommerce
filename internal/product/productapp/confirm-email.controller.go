package productapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

func (a *app) confirmEmailController(c *gin.Context) {
	confirmToken := c.Param("confirm_token")
	if confirmToken == "" {
		response.Send(c, a.log, nil, errs.New(errs.InvalidArgument, errors.New("missing confirm_token")))
	}
	err := a.confirmEmail(c.Request.Context(), confirmToken)
	response.Send(c, a.log, nil, err)
}

func (a *app) confirmEmail(ctx context.Context, confirmToken string) error {
	err := a.userBus.ConfirmEmail(ctx, confirmToken)
	if err != nil {
		if errors.Is(err, userbus.ErrNotFound) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "confirmEmail: confirmToken[%s]: %s", confirmToken, err)
	}

	return nil
}
