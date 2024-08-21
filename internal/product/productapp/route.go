package productapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authenticate := mid.Authenticate(a.log, a.auth)
	roleAdmin := mid.Authorize(a.log, a.auth, auth.Rules.Admin)

	r.GET("/products", a.queryController)
	r.GET("/products/:product_id", a.queryByIDController)
	r.POST("/products", authenticate, roleAdmin, a.createController)
	r.PUT("/products/:product_id", authenticate, roleAdmin, a.updateController)
	r.DELETE("/products/:product_id", authenticate, roleAdmin, a.deleteController)
}
