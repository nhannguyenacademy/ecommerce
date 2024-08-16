package orderbus

import (
	"context"
	"fmt"
)

func (b *Business) UpdateOrderStatus(ctx context.Context, ord Order, status Status) (Order, error) {
	if ord.Status.Equal(status) {
		return ord, nil
	}

	if ord.Status.Equal(Statuses.Finished) {
		return ord, fmt.Errorf("order %s: %w", ord.ID, ErrOrderAlreadyFinished)
	}

	if err := b.storer.UpdateStatus(ctx, ord, status); err != nil {
		return Order{}, fmt.Errorf("update status: %w", err)
	}

	return ord, nil
}
