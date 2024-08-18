package productbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (b *Business) Update(ctx context.Context, id uuid.UUID, input UpdateProduct) (Product, error) {
	prd, err := b.QueryByID(ctx, id)
	if err != nil {
		return Product{}, err
	}

	if input.Name != nil {
		prd.Name = *input.Name
	}

	if input.Description != nil {
		prd.Description = *input.Description
	}

	if input.ImageURL != nil {
		prd.ImageURL = *input.ImageURL
	}

	if input.Price != nil {
		prd.Price = *input.Price
	}

	if input.Quantity != nil {
		prd.Quantity = *input.Quantity
	}

	prd.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}
