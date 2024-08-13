package userapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authenticate := mid.Authenticate(a.log, a.auth)
	authorizeUser := mid.AuthorizeUser(a.log, a.auth, a.userBus, auth.Rules.User)

	r.POST("/users/register", a.registerController)
	r.POST("/users/login", a.loginController)
	r.GET("/users/confirm-email/:confirm_token", a.confirmEmailController)
	r.PUT("/users/:user_id", authenticate, authorizeUser, a.updateController)
	r.GET("/users/:user_id", authenticate, authorizeUser, a.queryByIDController)
}
