package productbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) Delete(ctx context.Context, id uuid.UUID) error {
	prd, err := b.QueryByID(ctx, id)
	if err != nil {
		return err
	}

	if err := b.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
