// Package userapp maintains the app layer api for the user domain.
package userapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"net/mail"
	"time"
)

type app struct {
	log       *logger.Logger
	auth      *auth.Auth
	activeKID string
	userBus   *userbus.Business
}

func New(
	log *logger.Logger,
	auth *auth.Auth,
	activeKID string,
	userBus *userbus.Business,
) *app {
	return &app{
		log:       log,
		auth:      auth,
		activeKID: activeKID,
		userBus:   userBus,
	}
}

func (a *app) registerHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	newUser, err := toBusRegisterUser(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	usr, err := a.userBus.Create(ctx, newUser)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			respond.Error(c, a.log, errs.New(errs.Aborted, userbus.ErrUniqueEmail))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "register: usr[%+v]: %s", usr, err))
		}
		return
	}

	// todo: send email confirmation, using queue
	// todo: redesign email confirmation flow, db
}

func (a *app) loginHandler(c *gin.Context) {
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

func (a *app) confirmEmailHandler(c *gin.Context) {
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

func (a *app) updateHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "user missing in context: %s", err))
		return
	}

	updateUser, err := toBusUpdateUser(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	updatedUser, err := a.userBus.Update(ctx, usr, updateUser)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "update: userID[%s] req[%+v]: %s", usr.ID, req, err))
		return
	}

	respond.Success(c, a.log, toAppUser(updatedUser))
}

func (a *app) queryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	usr, err := mid.GetUser(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "querybyid: %s", err))
		return
	}

	respond.Success(c, a.log, toAppUser(usr))
}
