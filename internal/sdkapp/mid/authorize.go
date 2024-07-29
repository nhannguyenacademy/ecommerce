package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

func Authorize(l *logger.Logger, auth *auth.Auth, rule string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userID, err := GetUserID(ctx)
		if err != nil {
			response.Send(c, l, nil, errs.New(errs.Unauthenticated, err))
			return
		}

		claims := GetClaims(ctx)

		if err := auth.Authorize(ctx, claims, userID, rule); err != nil {
			response.Send(c, l, nil, errs.New(errs.Unauthenticated, err))
			return
		}

		c.Next()
	}
}
