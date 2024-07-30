package mid

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

func Authorize(l *logger.Logger, auth *auth.Auth, rule auth.Rule) gin.HandlerFunc {
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

// AuthorizeUser executes the specified role and extracts the specified
// user from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(l *logger.Logger, auth *auth.Auth, userBus *userbus.Business, rule auth.Rule) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var userID uuid.UUID

		id := c.Param("user_id")
		if id != "" {
			var err error
			userID, err = uuid.Parse(id)
			if err != nil {
				response.Send(c, l, nil, errs.New(errs.Unauthenticated, ErrInvalidID))
				return
			}

			usr, err := userBus.QueryByID(ctx, userID)
			if err != nil {
				switch {
				case errors.Is(err, userbus.ErrNotFound):
					response.Send(c, l, nil, errs.New(errs.Unauthenticated, err))
					return
				default:
					response.Send(c, l, nil, errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err))
					return
				}
			}

			ctx = setUser(ctx, usr)
		}

		claims := GetClaims(ctx)
		if err := auth.Authorize(ctx, claims, userID, rule); err != nil {
			response.Send(c, l, nil, errs.New(errs.Unauthenticated, err))
			return
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
