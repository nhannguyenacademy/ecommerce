package orderbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) QueryByID(ctx context.Context, ordID uuid.UUID) (Order, error) {
	ord, err := b.storer.QueryByID(ctx, ordID)
	if err != nil {
		return Order{}, fmt.Errorf("query: ordID[%s]: %w", ordID, err)
	}

	return ord, nil
}
