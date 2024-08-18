package userbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Create adds a new user to the system.
func (b *Business) Create(ctx context.Context, nu NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
		ID:                uuid.New(),
		Name:              nu.Name,
		Email:             nu.Email,
		PasswordHash:      string(hash),
		Roles:             nu.Roles,
		Enabled:           true,
		EmailConfirmToken: nu.EmailConfirmToken,
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
