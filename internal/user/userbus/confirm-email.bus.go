package userbus

import (
	"context"
	"fmt"
)

// ConfirmEmail ...
func (b *Business) ConfirmEmail(ctx context.Context, confirmToken string) (User, error) {
	usr, err := b.storer.QueryByEmailConfirmToken(ctx, confirmToken)
	if err != nil {
		return User{}, fmt.Errorf("query: confirmToken[%s]: %w", confirmToken, err)
	}

	usr, err = b.Update(ctx, usr, UpdateUser{EmailConfirmToken: new(string)})
	if err != nil {
		return User{}, fmt.Errorf("update: confirmToken[%s]: %w", confirmToken, err)
	}

	return usr, nil
}
