package orderbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Create adds a new product to the system.
func (b *Business) Create(ctx context.Context, np NewProduct) (Product, error) {
	now := time.Now()

	prd := Product{
		ID:          uuid.New(),
		Name:        np.Name,
		Description: np.Description,
		ImageURL:    np.ImageURL,
		Price:       np.Price,
		Quantity:    np.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}
