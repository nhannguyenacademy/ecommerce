package productbus

import (
	"context"
	"fmt"
)

// Delete removes the specified product.
func (b *Business) Delete(ctx context.Context, prd Product) error {
	if err := b.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
