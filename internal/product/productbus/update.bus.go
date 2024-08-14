package productbus

import (
	"context"
	"fmt"
	"time"
)

// Update modifies information about a product.
func (b *Business) Update(ctx context.Context, prd Product, uu UpdateProduct) (Product, error) {
	if uu.Name != nil {
		prd.Name = *uu.Name
	}

	if uu.Description != nil {
		prd.Description = *uu.Description
	}

	if uu.ImageURL != nil {
		prd.ImageURL = *uu.ImageURL
	}

	if uu.Price != nil {
		prd.Price = *uu.Price
	}

	if uu.Quantity != nil {
		prd.Quantity = *uu.Quantity
	}

	prd.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}
