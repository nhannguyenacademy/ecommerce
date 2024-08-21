package userapp

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"net/mail"
	"time"
)

// =============================================================================

type user struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	DateCreated string   `json:"date_created"`
	DateUpdated string   `json:"date_updated"`
}

func toAppUser(bus userbus.User) user {
	return user{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		Email:       bus.Email.Address,
		Roles:       userbus.ParseRolesToString(bus.Roles),
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

// =============================================================================

type authenUser struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// =============================================================================

type registerReq struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"eqfield=Password"`
}

func toBusRegisterUser(app registerReq) (userbus.NewUser, error) {
	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	name, err := userbus.ParseName(app.Name)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := userbus.NewUser{
		Name:              name,
		Email:             *addr,
		Password:          app.Password,
		Roles:             []userbus.Role{userbus.Roles.User},
		EmailConfirmToken: uuid.NewString(),
	}

	return bus, nil
}

// =============================================================================

// loginUser defines the data needed to login a user.
type loginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// =============================================================================

// updateUserReq defines the data needed to update a user.
type updateUserReq struct {
	Name            *string `json:"name"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"password_confirm" binding:"omitempty,eqfield=Password"`
}

func toBusUpdateUser(app updateUserReq) (userbus.UpdateUser, error) {
	var name *userbus.Name
	if app.Name != nil {
		nm, err := userbus.ParseName(*app.Name)
		if err != nil {
			return userbus.UpdateUser{}, fmt.Errorf("parse: %w", err)
		}
		name = &nm
	}

	bus := userbus.UpdateUser{
		Name:     name,
		Password: app.Password,
	}

	return bus, nil
}
