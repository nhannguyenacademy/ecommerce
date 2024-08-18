// Package orderapp maintains the app layer.
package orderapp

import (
	"context"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

type app struct {
	log        *logger.Logger
	auth       *auth.Auth
	orderBus   *orderbus.Business
	dbBeginner sqldb.Beginner
}

func New(
	log *logger.Logger,
	auth *auth.Auth,
	dbBeginner sqldb.Beginner,
	orderBus *orderbus.Business,
) *app {
	return &app{
		log:        log,
		auth:       auth,
		dbBeginner: dbBeginner,
		orderBus:   orderBus,
	}
}

// newWithTx constructs a new app value using a store transaction that was created via middleware.
func (a *app) newWithTx(ctx context.Context) (*app, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	orderBus, err := a.orderBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	app := app{
		log:      a.log,
		auth:     a.auth,
		orderBus: orderBus,
	}

	return &app, nil
}
