package productbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// QueryByID finds the user by the specified ID.
func (b *Business) QueryByID(ctx context.Context, prdID uuid.UUID) (Product, error) {
	user, err := b.storer.QueryByID(ctx, prdID)
	if err != nil {
		return Product{}, fmt.Errorf("query: prdID[%s]: %w", prdID, err)
	}

	return user, nil
}
