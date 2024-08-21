// Package orderbus provides business access to product domain.
package orderbus

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"time"
)

var (
	ErrNotFound              = errors.New("order not found")
	ErrOrderAlreadyFinished  = errors.New("order already finished")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
)

type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, ord Order) error
	UpdateStatus(ctx context.Context, ord Order, status Status) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Order, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, ordID uuid.UUID) (Order, error)
	Delete(ctx context.Context, ord Order) error

	QueryOrderItems(ctx context.Context, ord Order) ([]OrderItem, error)
	DeleteOrderItems(ctx context.Context, ord Order) error
	CreateOrderItems(ctx context.Context, itms []OrderItem) error
}

// Business manages the set of APIs for user access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a business API for use.
func NewBusiness(log *logger.Logger, orderStorer Storer) *Business {
	return &Business{
		log:    log,
		storer: orderStorer,
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

func (b *Business) Delete(ctx context.Context, ord Order) error {
	if ord.Status.Equal(Statuses.Finished) {
		return fmt.Errorf("order %s: %w", ord.ID, ErrOrderAlreadyFinished)
	}

	if err := b.storer.Delete(ctx, ord); err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	if err := b.storer.DeleteOrderItems(ctx, ord); err != nil {
		return fmt.Errorf("delete order items: %w", err)
	}

	return nil
}

func (b *Business) QueryByID(ctx context.Context, id uuid.UUID) (Order, error) {
	ord, err := b.storer.QueryByID(ctx, id)
	if err != nil {
		return Order{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	return ord, nil
}

func (b *Business) QueryByIDWithItems(ctx context.Context, id uuid.UUID) (OrderWithItems, error) {
	ord, err := b.storer.QueryByID(ctx, id)
	if err != nil {
		return OrderWithItems{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	itms, err := b.storer.QueryOrderItems(ctx, ord)
	if err != nil {
		return OrderWithItems{}, fmt.Errorf("query order items: %w", err)
	}

	return OrderWithItems{
		Order: ord,
		Items: itms,
	}, nil
}

func (b *Business) UpdateStatus(ctx context.Context, ord Order, status Status) (Order, error) {
	if ord.Status.Equal(status) {
		return ord, nil
	}

	if ord.Status.Equal(Statuses.Finished) {
		return ord, fmt.Errorf("order %s: %w", ord.ID, ErrOrderAlreadyFinished)
	}

	if ord.Status.Equal(Statuses.Cancelled) {
		return ord, fmt.Errorf("order %s: %w", ord.ID, ErrOrderAlreadyCancelled)
	}

	if err := b.storer.UpdateStatus(ctx, ord, status); err != nil {
		return Order{}, fmt.Errorf("update status: %w", err)
	}

	return ord, nil
}

func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Order, error) {
	ords, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return ords, nil
}

func (b *Business) Create(ctx context.Context, input NewOrder) (Order, error) {
	var (
		orderAmount int64
		ordID       = uuid.New()
		ordItms     = make([]OrderItem, len(input.Items))
		now         = time.Now()
	)

	for i, itm := range input.Items {
		ordItms[i] = OrderItem{
			ID:          uuid.New(),
			OrderID:     ordID,
			ProductID:   itm.ProductID,
			Quantity:    itm.Quantity,
			Price:       itm.Price,
			DateCreated: now,
			DateUpdated: now,
		}

		orderAmount += itm.Price
	}

	ord := Order{
		ID:          ordID,
		UserID:      input.UserID,
		Amount:      orderAmount,
		Status:      Statuses.Created,
		DateCreated: now,
		DateUpdated: now,
	}
	if err := b.storer.Create(ctx, ord); err != nil {
		return Order{}, fmt.Errorf("create: %w", err)
	}

	if err := b.storer.CreateOrderItems(ctx, ordItms); err != nil {
		return Order{}, fmt.Errorf("create order items: %w", err)
	}

	return ord, nil
}
