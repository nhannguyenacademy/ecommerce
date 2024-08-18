package orderapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authenticate := mid.Authenticate(a.log, a.auth)
	authorizeAdmin := mid.Authorize(a.log, a.auth, auth.Rules.Admin)
	transaction := mid.BeginCommitRollback(a.log, a.dbBeginner)

	r.POST("/orders", authenticate, authorizeAdmin, transaction, a.createController)
}
