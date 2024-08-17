// Package orderbus provides business access to product domain.
package orderbus

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

// Set of error variables for CRUD operations.

var (
	ErrNotFound             = errors.New("order not found")
	ErrMissingProducts      = errors.New("missing products")
	ErrOrderAlreadyFinished = errors.New("order already finished")
)

type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, ord Order) error
	UpdateStatus(ctx context.Context, ord Order, status Status) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Order, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, ordID uuid.UUID) (Order, error)
	Delete(ctx context.Context, ord Order) error
}

// Business manages the set of APIs for user access.
type Business struct {
	log        *logger.Logger
	storer     Storer
	productBus *productbus.Business
}

// NewBusiness constructs a business API for use.
func NewBusiness(
	log *logger.Logger,
	orderStorer Storer,
	productBus *productbus.Business,
) *Business {
	return &Business{
		log:        log,
		storer:     orderStorer,
		productBus: productBus,
	}
}

func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (*Business, error) {
	orderStorer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:    b.log,
		storer: orderStorer,
	}

	return &bus, nil
}
