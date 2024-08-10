package userapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"net/mail"
	"time"
)

func (a *app) loginController(c *gin.Context) {
	var lu loginUser
	if err := c.ShouldBindJSON(&lu); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
		}

		response.Send(c, a.log, nil, err)
		return
	}

	usr, err := a.login(c.Request.Context(), lu)
	response.Send(c, a.log, usr, err)
}

func (a *app) login(ctx context.Context, lu loginUser) (authenUser, error) {
	addr, err := mail.ParseAddress(lu.Email)
	if err != nil {
		return authenUser{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Authenticate(ctx, *addr, lu.Password)
	if err != nil {
		if errors.Is(err, userbus.ErrAuthenticationFailure) {
			return authenUser{}, errs.New(errs.Unauthenticated, err)
		}
		return authenUser{}, errs.New(errs.Internal, err)
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: userbus.ParseRolesToString(usr.Roles),
	}
	token, err := a.auth.GenerateToken(a.activeKID, claims)
	if err != nil {
		return authenUser{}, errs.New(errs.Internal, err)
	}

	return authenUser{
		User:  toAppUser(usr),
		Token: token,
	}, nil
}
