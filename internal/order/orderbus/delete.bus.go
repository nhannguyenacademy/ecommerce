package orderbus

import (
	"context"
	"fmt"
)

func (b *Business) Delete(ctx context.Context, ord Order) error {
	// todo: check if order has any success payments,
	// but orderbus cannot import paymentbus, use delegate instead

	if ord.Status.Equal(Statuses.Finished) {
		return fmt.Errorf("order %s: %w", ord.ID, ErrOrderAlreadyFinished)
	}

	if err := b.storer.Delete(ctx, ord); err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	return nil
}
