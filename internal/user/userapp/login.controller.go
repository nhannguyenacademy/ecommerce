package userapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"net/mail"
	"time"
)

func (a *app) loginController(c *gin.Context) {
	ctx := c.Request.Context()

	var req loginUser
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	usr, err := a.userBus.Authenticate(ctx, *addr, req.Password)
	if err != nil {
		if errors.Is(err, userbus.ErrAuthenticationFailure) {
			respond.Error(c, a.log, errs.New(errs.Unauthenticated, err))
		} else {
			respond.Error(c, a.log, errs.New(errs.Internal, err))
		}
		return
	}

	if !usr.Enabled || usr.EmailConfirmToken != "" {
		respond.Error(c, a.log, errs.New(errs.Unauthenticated, errors.New("invalid user")))
		return
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(365 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: userbus.ParseRolesToString(usr.Roles),
	}
	token, err := a.auth.GenerateToken(a.activeKID, claims)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.Internal, err))
		return
	}

	respond.Success(c, a.log, authenUser{
		UserID: usr.ID.String(),
		Token:  token,
	})
}
