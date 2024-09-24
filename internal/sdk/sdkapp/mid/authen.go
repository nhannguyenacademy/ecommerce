package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

func Authenticate(l *logger.Logger, auth *auth.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		claims, err := auth.Authenticate(ctx, c.GetHeader("Authorization"))
		if err != nil {
			respond.Error(c, l, errs.New(errs.Unauthenticated, err))
			return
		}

		if claims.Subject == "" {
			respond.Error(c, l, errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims"))
			return
		}

		subjectID, err := uuid.Parse(claims.Subject)
		if err != nil {
			respond.Error(c, l, errs.Newf(errs.Unauthenticated, "parsing subject: %s", err))
			return
		}

		ctx = setUserID(ctx, subjectID)
		ctx = setClaims(ctx, claims)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
