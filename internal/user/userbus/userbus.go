// Package userbus provides business access to user domain.
package userbus

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to perists and retrieve data.
type Storer interface {
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
	QueryByEmailConfirmToken(ctx context.Context, token string) (User, error)
}

// Business manages the set of APIs for user access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a user business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

func (b *Business) Create(ctx context.Context, newUser NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	user := User{
		ID:                uuid.New(),
		Name:              newUser.Name,
		Email:             newUser.Email,
		PasswordHash:      string(hash),
		Roles:             newUser.Roles,
		Enabled:           true,
		EmailConfirmToken: newUser.EmailConfirmToken,
		DateCreated:       now,
		DateUpdated:       now,
	}

	if err := b.storer.Create(ctx, user); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return user, nil
}

func (b *Business) Update(ctx context.Context, user User, updateUser UpdateUser) (User, error) {
	if updateUser.Name != nil {
		user.Name = *updateUser.Name
	}

	if updateUser.Email != nil {
		user.Email = *updateUser.Email
	}

	if updateUser.Roles != nil {
		user.Roles = updateUser.Roles
	}

	if updateUser.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
		}
		user.PasswordHash = string(pw)
	}

	if updateUser.Enabled != nil {
		user.Enabled = *updateUser.Enabled
	}

	if updateUser.EmailConfirmToken != nil {
		user.EmailConfirmToken = *updateUser.EmailConfirmToken
	}

	user.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, user); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	return user, nil
}

func (b *Business) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := b.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
	}

	return user, nil
}

func (b *Business) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {
	user, err := b.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, fmt.Errorf("query: email[%s]: %w", email, ErrAuthenticationFailure)
		}
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return user, nil
}

func (b *Business) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	user, err := b.storer.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	return user, nil
}

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
