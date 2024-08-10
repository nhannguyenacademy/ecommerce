// Package userapp maintains the app layer api for the user domain.
package userapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
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

func (a *app) Routes(r gin.IRouter) {
	authen := mid.Authenticate(a.log, a.auth)
	ruleAdmin := mid.Authorize(a.log, a.auth, auth.Rules.Admin)
	ruleAuthorizeUser := mid.AuthorizeUser(a.log, a.auth, a.userBus, auth.Rules.AdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(a.log, a.auth, a.userBus, auth.Rules.Admin)

	r.POST("/users/register", a.registerController)
	r.POST("/users/login", a.loginController)

	r.GET("/users", authen, ruleAdmin, a.queryController)
	r.GET("/users/:user_id", authen, ruleAuthorizeUser, a.queryByIDController)
	r.POST("/users", authen, ruleAdmin, a.createController)
	r.PUT("/users/:user_id", authen, ruleAuthorizeUser, a.updateController)
	r.DELETE("/users/:user_id", authen, ruleAuthorizeAdmin, a.deleteController)
}
