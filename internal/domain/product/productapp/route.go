package productapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authenticate := mid.Authenticate(a.log, a.auth)
	roleAdmin := mid.Authorize(a.log, a.auth, auth.Rules.Admin)

	r.GET("/products", a.queryHandler)
	r.GET("/products/:product_id", a.queryByIDHandler)
	r.POST("/products", authenticate, roleAdmin, a.createHandler)
	r.PUT("/products/:product_id", authenticate, roleAdmin, a.updateHandler)
	r.DELETE("/products/:product_id", authenticate, roleAdmin, a.deleteHandler)
}
