package userapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log     *logger.Logger
	UserBus *userbus.Business
	Auth    *auth.Auth
}

func Routes(r gin.IRouter, cfg Config) {
	authen := mid.Authenticate(cfg.Log, cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Log, cfg.Auth, auth.Rules.Admin)
	ruleAuthorizeUser := mid.AuthorizeUser(cfg.Log, cfg.Auth, cfg.UserBus, auth.Rules.AdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(cfg.Log, cfg.Auth, cfg.UserBus, auth.Rules.Admin)

	app := NewApp(cfg.UserBus)

	r.GET("/users", authen, ruleAdmin, func(c *gin.Context) {
		results, err := app.query(c.Request.Context(), parseQueryParams(c.Request))
		response.Send(c, cfg.Log, results, err)
	})

	r.GET("/users/:user_id", authen, ruleAuthorizeUser, func(c *gin.Context) {
		u, err := app.queryByID(c.Request.Context())
		response.Send(c, cfg.Log, u, err)
	})

	r.POST("/users", authen, ruleAdmin, func(c *gin.Context) {
		var nu NewUser
		if err := c.ShouldBindJSON(&nu); err != nil {
			var vErrs validator.ValidationErrors
			if errors.As(err, &vErrs) {
				err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
			}

			response.Send(c, cfg.Log, nil, err)
			return
		}

		u, err := app.create(c.Request.Context(), nu)
		response.Send(c, cfg.Log, u, err)
	})

	r.PUT("/users/role/:user_id", authen, ruleAuthorizeAdmin, func(c *gin.Context) {
		var ur UpdateUserRole
		if err := c.ShouldBindJSON(&ur); err != nil {
			var vErrs validator.ValidationErrors
			if errors.As(err, &vErrs) {
				err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
			}

			response.Send(c, cfg.Log, nil, err)
			return
		}

		u, err := app.updateRole(c.Request.Context(), ur)
		response.Send(c, cfg.Log, u, err)
	})

	r.PUT("/users/:user_id", authen, ruleAuthorizeUser, func(c *gin.Context) {
		var uu UpdateUser
		if err := c.ShouldBindJSON(&uu); err != nil {
			var vErrs validator.ValidationErrors
			if errors.As(err, &vErrs) {
				err = errs.Newf(errs.InvalidArgument, "%s", vErrs)
			}

			response.Send(c, cfg.Log, nil, err)
			return
		}

		u, err := app.update(c.Request.Context(), uu)
		response.Send(c, cfg.Log, u, err)
	})

	r.DELETE("/users/:user_id", authen, ruleAuthorizeAdmin, func(c *gin.Context) {
		err := app.delete(c.Request.Context())
		response.Send(c, cfg.Log, nil, err)
	})
}
