package userbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// QueryByID finds the user by the specified Ib.
func (b *Business) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := b.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
	}

	return user, nil
}