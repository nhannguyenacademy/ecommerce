package orderbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (b *Business) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status Status) (Order, error) {
	ord, err := b.QueryByID(ctx, id)
	if err != nil {
		return Order{}, err
	}

	if ord.Status.Equal(status) {
		return ord, nil
	}

	if ord.Status.Equal(Statuses.Finished) {
		return ord, fmt.Errorf("order %s: %w", ord.ID, ErrOrderAlreadyFinished)
	}

	if ord.Status.Equal(Statuses.Cancelled) {
		return ord, fmt.Errorf("order %s: %w", ord.ID, ErrOrderAlreadyCancelled)
	}

	if err := b.storer.UpdateStatus(ctx, ord, status); err != nil {
		return Order{}, fmt.Errorf("update status: %w", err)
	}

	return ord, nil
}
