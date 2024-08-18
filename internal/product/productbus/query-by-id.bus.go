package productbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) QueryByID(ctx context.Context, id uuid.UUID) (Product, error) {
	prd, err := b.storer.QueryByID(ctx, id)
	if err != nil {
		return Product{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	return prd, nil
}
