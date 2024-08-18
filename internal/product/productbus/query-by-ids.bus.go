package productbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) QueryByIDs(ctx context.Context, ids []uuid.UUID) ([]Product, error) {
	prds, err := b.storer.QueryByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("query: ids[%+v]: %w", ids, err)
	}

	return prds, nil
}
