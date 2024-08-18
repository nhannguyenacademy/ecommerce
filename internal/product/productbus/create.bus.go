package productbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (b *Business) Create(ctx context.Context, input NewProduct) (Product, error) {
	now := time.Now()

	// todo: upload image to s3

	prd := Product{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
		Quantity:    input.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}
