package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

func Authen(l *logger.Logger, auth *auth.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		claims, err := auth.Authenticate(ctx, c.GetHeader("Authorization"))
		if err != nil {
			response.Send(c, l, nil, errs.New(errs.Unauthenticated, err))
			return
		}

		if claims.Subject == "" {
			response.Send(c, l, nil, errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims"))
			return
		}

		subjectID, err := uuid.Parse(claims.Subject)
		if err != nil {
			response.Send(c, l, nil, errs.Newf(errs.Unauthenticated, "parsing subject: %s", err))
			return
		}

		ctx = setUserID(ctx, subjectID)
		ctx = setClaims(ctx, claims)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
