package userapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authen := mid.Authenticate(a.log, a.auth)
	ruleAdmin := mid.Authorize(a.log, a.auth, auth.Rules.Admin)
	ruleAuthorizeUser := mid.AuthorizeUser(a.log, a.auth, a.userBus, auth.Rules.AdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(a.log, a.auth, a.userBus, auth.Rules.Admin)

	// Guests routes
	r.POST("/users/register", a.registerController)
	r.POST("/users/login", a.loginController)
	r.GET("/users/confirm-email/:confirm_token", a.confirmEmailController)

	// Users or admin routes
	r.GET("/users/:user_id", authen, ruleAuthorizeUser, a.queryByIDController)
	r.PUT("/users/:user_id", authen, ruleAuthorizeUser, a.updateController)
	r.GET("/users", authen, ruleAdmin, a.queryController)
	r.POST("/users", authen, ruleAdmin, a.createController)
	r.DELETE("/users/:user_id", authen, ruleAuthorizeAdmin, a.deleteController)
}
