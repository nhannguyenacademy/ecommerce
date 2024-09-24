package orderapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authenticate := mid.Authenticate(a.log, a.auth)
	roleAdmin := mid.Authorize(a.log, a.auth, auth.Rules.Admin)
	orderOwner := mid.AuthorizeOrder(a.log, a.auth, a.orderBus, auth.Rules.Owner)
	adminOrOrderOwner := mid.AuthorizeOrder(a.log, a.auth, a.orderBus, auth.Rules.AdminOrOwner)
	transaction := mid.BeginCommitRollback(a.log, a.dbBeginner)

	r.POST("/orders", authenticate, transaction, a.createHandler)
	r.PUT("/orders/:order_id/cancel", authenticate, orderOwner, a.cancelHandler)
	r.GET("/orders/:order_id", authenticate, adminOrOrderOwner, a.queryByIDHandler)
	r.GET("/:user_id/orders", authenticate, a.queryUserOrdersHandler)
	r.GET("/orders", authenticate, roleAdmin, a.queryHandler)
	r.PUT("/orders/:order_id", authenticate, roleAdmin, a.updateStatusHandler)
	r.DELETE("/orders/:order_id", authenticate, roleAdmin, transaction, a.deleteHandler)
}
