package userapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log     *logger.Logger
	UserBus *userbus.Business
	//AuthClient *authclient.Client
}

func Routes(r gin.IRouter, cfg Config) {
	app := NewApp(cfg.UserBus)

	r.GET("/users", func(c *gin.Context) {
		results, err := app.query(c.Request.Context(), parseQueryParams(c.Request))
		response.Send(c, cfg.Log, results, err)
	})

	//r.POST("/users", func(c *gin.Context) {
	//	results, err := app.create(c.Request.Context())
	//	response.Send(c, cfg.Log, results, err)
	//})
}
