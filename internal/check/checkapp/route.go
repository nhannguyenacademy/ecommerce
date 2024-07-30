package checkapp

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

func Routes(r gin.IRouter, cfg Config) {
	app := NewApp(cfg.Build, cfg.Log, cfg.DB)

	r.GET("/readiness", func(c *gin.Context) {
		err := app.readiness(c.Request.Context())
		response.Send(c, cfg.Log, nil, err)
	})

	r.GET("/liveness", func(c *gin.Context) {
		info := app.liveness(c.Request.Context())
		response.Send(c, cfg.Log, info, nil)
	})
}
