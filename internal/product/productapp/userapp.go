// Package productapp maintains the app layer api for the product domain.
package productapp

import (
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

type app struct {
	log        *logger.Logger
	auth       *auth.Auth
	productBus *productbus.Business
}

func New(
	log *logger.Logger,
	auth *auth.Auth,
	productBus *productbus.Business,
) *app {
	return &app{
		log:        log,
		auth:       auth,
		productBus: productBus,
	}
}
