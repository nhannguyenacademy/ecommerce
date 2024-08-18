package userbus

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (b *Business) Update(ctx context.Context, usr User, input UpdateUser) (User, error) {
	if input.Name != nil {
		usr.Name = *input.Name
	}

	if input.Email != nil {
		usr.Email = *input.Email
	}

	if input.Roles != nil {
		usr.Roles = input.Roles
	}

	if input.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
		}
		usr.PasswordHash = string(pw)
	}

	if input.Enabled != nil {
		usr.Enabled = *input.Enabled
	}

	if input.EmailConfirmToken != nil {
		usr.EmailConfirmToken = *input.EmailConfirmToken
	}

	usr.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	return usr, nil
}
