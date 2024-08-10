package userbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Register a new user to the system.
func (b *Business) Register(ctx context.Context, nu RegisterUser) (User, error) {
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
		Roles:             []Role{Roles.User},
		Enabled:           true,
		EmailConfirmToken: uuid.NewString(),
		DateCreated:       now,
		DateUpdated:       now,
	}

	if err := b.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}
