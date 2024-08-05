package userbus

import (
	"context"
	"fmt"
)

// Delete removes the specified user.
func (b *Business) Delete(ctx context.Context, usr User) error {
	if err := b.storer.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
