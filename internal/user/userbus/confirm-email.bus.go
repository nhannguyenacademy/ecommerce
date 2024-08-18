package userbus

import (
	"context"
	"fmt"
)

func (b *Business) ConfirmEmail(ctx context.Context, confirmToken string) error {
	usr, err := b.storer.QueryByEmailConfirmToken(ctx, confirmToken)
	if err != nil {
		return fmt.Errorf("query: confirmToken[%s]: %w", confirmToken, err)
	}

	_, err = b.Update(ctx, usr, UpdateUser{EmailConfirmToken: new(string)})
	if err != nil {
		return fmt.Errorf("update: confirmToken[%s]: %w", confirmToken, err)
	}

	return nil
}
