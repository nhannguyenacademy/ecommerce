// Package orderbus provides business access to product domain.
package orderbus

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
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
	Create(ctx context.Context, order Order) error
	UpdateStatus(ctx context.Context, order Order, status Status) error
	Query(ctx context.Context, filter QueryFilter, sortBy sort.By, page page.Page) ([]Order, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, orderID uuid.UUID) (Order, error)
	Delete(ctx context.Context, order Order) error

	QueryOrderItems(ctx context.Context, order Order) ([]OrderItem, error)
	DeleteOrderItems(ctx context.Context, order Order) error
	CreateOrderItems(ctx context.Context, items []OrderItem) error
}

// Business manages the set of APIs for user access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (*Business, error) {
	storerTx, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:    b.log,
		storer: storerTx,
	}

	return &bus, nil
}

func (b *Business) Delete(ctx context.Context, order Order) error {
	if order.Status.Equal(Statuses.Finished) {
		return fmt.Errorf("order %s: %w", order.ID, ErrOrderAlreadyFinished)
	}

	if err := b.storer.Delete(ctx, order); err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	if err := b.storer.DeleteOrderItems(ctx, order); err != nil {
		return fmt.Errorf("delete order items: %w", err)
	}

	return nil
}

func (b *Business) QueryByID(ctx context.Context, id uuid.UUID) (Order, error) {
	order, err := b.storer.QueryByID(ctx, id)
	if err != nil {
		return Order{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	return order, nil
}

func (b *Business) QueryByIDWithItems(ctx context.Context, id uuid.UUID) (OrderWithItems, error) {
	order, err := b.storer.QueryByID(ctx, id)
	if err != nil {
		return OrderWithItems{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	items, err := b.storer.QueryOrderItems(ctx, order)
	if err != nil {
		return OrderWithItems{}, fmt.Errorf("query order items: %w", err)
	}

	return OrderWithItems{
		Order: order,
		Items: items,
	}, nil
}

func (b *Business) UpdateStatus(ctx context.Context, order Order, status Status) (Order, error) {
	if order.Status.Equal(status) {
		return order, nil
	}

	if order.Status.Equal(Statuses.Finished) {
		return order, fmt.Errorf("order %s: %w", order.ID, ErrOrderAlreadyFinished)
	}

	if order.Status.Equal(Statuses.Cancelled) {
		return order, fmt.Errorf("order %s: %w", order.ID, ErrOrderAlreadyCancelled)
	}

	if err := b.storer.UpdateStatus(ctx, order, status); err != nil {
		return Order{}, fmt.Errorf("update status: %w", err)
	}

	return order, nil
}

func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

func (b *Business) Query(ctx context.Context, filter QueryFilter, sortBy sort.By, page page.Page) ([]Order, error) {
	ords, err := b.storer.Query(ctx, filter, sortBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return ords, nil
}

func (b *Business) Create(ctx context.Context, newOrder NewOrder) (Order, error) {
	var (
		orderAmount int64
		orderID     = uuid.New()
		orderItems  = make([]OrderItem, len(newOrder.Items))
		now         = time.Now()
	)

	for i, item := range newOrder.Items {
		orderItems[i] = OrderItem{
			ID:          uuid.New(),
			OrderID:     orderID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			Price:       item.Price,
			DateCreated: now,
			DateUpdated: now,
		}

		orderAmount += item.Price
	}

	order := Order{
		ID:          orderID,
		UserID:      newOrder.UserID,
		Amount:      orderAmount,
		Status:      Statuses.Created,
		DateCreated: now,
		DateUpdated: now,
	}
	if err := b.storer.Create(ctx, order); err != nil {
		return Order{}, fmt.Errorf("create: %w", err)
	}

	if err := b.storer.CreateOrderItems(ctx, orderItems); err != nil {
		return Order{}, fmt.Errorf("create order items: %w", err)
	}

	return order, nil
}
