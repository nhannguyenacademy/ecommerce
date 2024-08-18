package userbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (b *Business) Create(ctx context.Context, input NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
		ID:                uuid.New(),
		Name:              input.Name,
		Email:             input.Email,
		PasswordHash:      string(hash),
		Roles:             input.Roles,
		Enabled:           true,
		EmailConfirmToken: input.EmailConfirmToken,
		DateCreated:       now,
		DateUpdated:       now,
	}

	if err := b.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	if usr.EmailConfirmToken != "" {
		// todo: send email confirmation
	}

	return usr, nil
}
