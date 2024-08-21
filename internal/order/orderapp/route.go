package orderapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authenticate := mid.Authenticate(a.log, a.auth)
	roleUser := mid.Authorize(a.log, a.auth, auth.Rules.User)
	roleAdmin := mid.Authorize(a.log, a.auth, auth.Rules.Admin)
	orderOwner := mid.AuthorizeOrder(a.log, a.auth, a.orderBus, auth.Rules.Owner)
	adminOrOrderOwner := mid.AuthorizeOrder(a.log, a.auth, a.orderBus, auth.Rules.AdminOrOwner)
	transaction := mid.BeginCommitRollback(a.log, a.dbBeginner)

	r.POST("/orders", authenticate, roleUser, transaction, a.createController)
	r.PUT("/orders/:order_id/cancel", authenticate, orderOwner, a.cancelController)
	r.GET("/orders/:order_id", authenticate, adminOrOrderOwner, a.queryByIDController)
	r.GET("/:user_id/orders", authenticate, adminOrOrderOwner)
	r.GET("/orders", authenticate, roleAdmin)
	r.PUT("/orders/:order_id", authenticate, roleAdmin, a.updateStatusController)
	r.DELETE("/orders/:order_id", authenticate, roleAdmin, transaction, a.deleteController)
}
