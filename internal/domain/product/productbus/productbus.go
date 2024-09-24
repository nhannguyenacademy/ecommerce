// Package productbus provides business access to product domain.
package productbus

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sort"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"time"
)

// Set of error variables for CRUD operations.

var (
	ErrNotFound = errors.New("product not found")
)

// Storer interface declares the behavior this package needs to perists and retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, product Product) error
	Update(ctx context.Context, product Product) error
	Delete(ctx context.Context, product Product) error
	Query(ctx context.Context, filter QueryFilter, sortBy sort.By, page page.Page) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByIDs(ctx context.Context, productIDs []uuid.UUID) ([]Product, error)
}

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

// NewWithTx constructs a new business value that will use the specified transaction in any store related calls.
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

func (b *Business) Create(ctx context.Context, newProduct NewProduct) (Product, error) {
	now := time.Now()

	// todo: upload image to s3

	product := Product{
		ID:          uuid.New(),
		Name:        newProduct.Name,
		Description: newProduct.Description,
		ImageURL:    newProduct.ImageURL,
		Price:       newProduct.Price,
		Quantity:    newProduct.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, product); err != nil {
		return Product{}, fmt.Errorf("create: %w", err)
	}

	return product, nil
}

func (b *Business) Update(ctx context.Context, product Product, updateProduct UpdateProduct) (Product, error) {
	if updateProduct.Name != nil {
		product.Name = *updateProduct.Name
	}

	if updateProduct.Description != nil {
		product.Description = *updateProduct.Description
	}

	if updateProduct.ImageURL != nil {
		product.ImageURL = *updateProduct.ImageURL
	}

	if updateProduct.Price != nil {
		product.Price = *updateProduct.Price
	}

	if updateProduct.Quantity != nil {
		product.Quantity = *updateProduct.Quantity
	}

	product.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, product); err != nil {
		return Product{}, fmt.Errorf("update: %w", err)
	}

	return product, nil
}

func (b *Business) Query(ctx context.Context, filter QueryFilter, sortBy sort.By, page page.Page) ([]Product, error) {
	products, err := b.storer.Query(ctx, filter, sortBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return products, nil
}

func (b *Business) Delete(ctx context.Context, product Product) error {
	if err := b.storer.Delete(ctx, product); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

func (b *Business) QueryByIDs(ctx context.Context, productIDs []uuid.UUID) ([]Product, error) {
	products, err := b.storer.QueryByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("query: productIDs[%+v]: %w", productIDs, err)
	}

	return products, nil
}

func (b *Business) QueryByID(ctx context.Context, productID uuid.UUID) (Product, error) {
	product, err := b.storer.QueryByID(ctx, productID)
	if err != nil {
		return Product{}, fmt.Errorf("query: productID[%s]: %w", productID, err)
	}

	return product, nil
}
