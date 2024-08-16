package productbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) QueryByID(ctx context.Context, prdID uuid.UUID) (Product, error) {
	prd, err := b.storer.QueryByID(ctx, prdID)
	if err != nil {
		return Product{}, fmt.Errorf("query: prdID[%s]: %w", prdID, err)
	}

	return prd, nil
}
