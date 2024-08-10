package userbus

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Update modifies information about a user.
func (b *Business) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
	if uu.Name != nil {
		usr.Name = *uu.Name
	}

	if uu.Email != nil {
		usr.Email = *uu.Email
	}

	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}

	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
		}
		usr.PasswordHash = string(pw)
	}

	if uu.Enabled != nil {
		usr.Enabled = *uu.Enabled
	}

	if uu.EmailConfirmToken != nil {
		usr.EmailConfirmToken = *uu.EmailConfirmToken
	}

	usr.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	return usr, nil
}
