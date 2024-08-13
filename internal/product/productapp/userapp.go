// Package userapp maintains the app layer api for the user domain.
package productapp

import (
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

type app struct {
	log       *logger.Logger
	auth      *auth.Auth
	activeKID string
	userBus   *userbus.Business
}

func New(
	log *logger.Logger,
	auth *auth.Auth,
	activeKID string,
	userBus *userbus.Business,
) *app {
	return &app{
		log:       log,
		auth:      auth,
		activeKID: activeKID,
		userBus:   userBus,
	}
}
