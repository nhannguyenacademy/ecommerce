package userapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"net/http"
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
		// app.query makes it easy to switching to new http framework
		results, err := app.query(c.Request.Context(), parseQueryParams(c.Request))
		if err != nil {
			appErr := errs.NewError(err)
			c.AbortWithError(appErr.HTTPStatus(), appErr)
			return
		}
		c.JSON(http.StatusOK, results)
	})
}
