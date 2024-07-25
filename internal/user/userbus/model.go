package userbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// User represents information about an individual user.
type User struct {
	ID           uuid.UUID
	Name         Name
	Email        mail.Address
	Roles        []Role
	PasswordHash []byte
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name     Name
	Email    mail.Address
	Roles    []Role
	Password string
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Name     *Name
	Email    *mail.Address
	Roles    []Role
	Password *string
	Enabled  *bool
}
