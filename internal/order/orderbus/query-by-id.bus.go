package orderbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) QueryByID(ctx context.Context, id uuid.UUID) (Order, error) {
	ord, err := b.storer.QueryByID(ctx, id)
	if err != nil {
		return Order{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	return ord, nil
}
