package productbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) QueryByIDs(ctx context.Context, prdIDs []uuid.UUID) ([]Product, error) {
	prds, err := b.storer.QueryByIDs(ctx, prdIDs)
	if err != nil {
		return nil, fmt.Errorf("query: prdIDs[%+v]: %w", prdIDs, err)
	}

	return prds, nil
}
