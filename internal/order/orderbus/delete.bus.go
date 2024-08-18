package orderbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) Delete(ctx context.Context, id uuid.UUID) error {
	ord, err := b.QueryByID(ctx, id)
	if err != nil {
		return err
	}

	// todo: check if order has any success payments, but orderbus cannot import paymentbus, use delegate instead

	if ord.Status.Equal(Statuses.Finished) {
		return fmt.Errorf("order %s: %w", ord.ID, ErrOrderAlreadyFinished)
	}

	if err := b.storer.Delete(ctx, ord); err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	return nil
}
