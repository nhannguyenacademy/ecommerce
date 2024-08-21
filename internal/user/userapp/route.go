package userapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authenticate := mid.Authenticate(a.log, a.auth)
	owner := mid.AuthorizeUser(a.log, a.auth, a.userBus, auth.Rules.Owner)

	r.POST("/users/register", a.registerHandler)
	r.POST("/users/login", a.loginHandler)
	r.GET("/users/confirm-email/:confirm_token", a.confirmEmailHandler)
	r.PUT("/users/:user_id", authenticate, owner, a.updateHandler)
	r.GET("/users/:user_id", authenticate, owner, a.queryByIDHandler)
}
