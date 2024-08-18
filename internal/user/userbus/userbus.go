// Package userbus provides business access to user domain.
package userbus

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"net/mail"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to perists and retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
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
func NewBusiness(
	log *logger.Logger,
	storer Storer,
) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}
