package orderapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
)

func (a *app) Routes(r gin.IRouter) {
	authenticate := mid.Authenticate(a.log, a.auth)
	authorizeAdmin := mid.Authorize(a.log, a.auth, auth.Rules.Admin)

	r.GET("/products", a.queryController)
	r.GET("/products/:product_id", a.queryByIDController)
	r.POST("/products", authenticate, authorizeAdmin, a.createController)
	r.PUT("/products/:product_id", authenticate, authorizeAdmin, a.updateController)
	r.DELETE("/products/:product_id", authenticate, authorizeAdmin, a.deleteController)
}
