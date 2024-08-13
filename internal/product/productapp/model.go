package productapp

import (
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"net/http"
	"net/mail"
	"time"
)

// =============================================================================
// Query params

// queryParams represents the set of possible query strings.
type queryParams struct {
	Page             string
	Rows             string
	OrderBy          string
	ID               string
	Name             string
	Email            string
	StartCreatedDate string
	EndCreatedDate   string
}

func parseQueryParams(r *http.Request) queryParams {
	values := r.URL.Query()

	filter := queryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("row"),
		OrderBy:          values.Get("orderBy"),
		ID:               values.Get("user_id"),
		Name:             values.Get("name"),
		Email:            values.Get("email"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
	}

	return filter
}

// =============================================================================

// user represents information about an individual user.
type user struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	Roles             []string `json:"roles"`
	PasswordHash      string   `json:"-"`
	Enabled           bool     `json:"-"`
	EmailConfirmToken string   `json:"-"`
	DateCreated       string   `json:"date_created"`
	DateUpdated       string   `json:"date_updated"`
}

func toAppUser(bus userbus.User) user {
	return user{
		ID:                bus.ID.String(),
		Name:              bus.Name.String(),
		Email:             bus.Email.Address,
		Roles:             userbus.ParseRolesToString(bus.Roles),
		PasswordHash:      bus.PasswordHash,
		Enabled:           bus.Enabled,
		EmailConfirmToken: bus.EmailConfirmToken,
		DateCreated:       bus.DateCreated.Format(time.RFC3339),
		DateUpdated:       bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []user {
	app := make([]user, len(users))
	for i, usr := range users {
		app[i] = toAppUser(usr)
	}

	return app
}

// =============================================================================

type authenUser struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// =============================================================================

// registerUser defines the data needed to register a new user.
type registerUser struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"eqfield=Password"`
}

func toBusRegisterUser(app registerUser) (userbus.RegisterUser, error) {
	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return userbus.RegisterUser{}, fmt.Errorf("parse: %w", err)
	}

	name, err := userbus.ParseName(app.Name)
	if err != nil {
		return userbus.RegisterUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := userbus.RegisterUser{
		Name:     name,
		Email:    *addr,
		Password: app.Password,
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

// newUser defines the data needed to add a new user.
type newUser struct {
	Name            string   `json:"name" binding:"required"`
	Email           string   `json:"email" binding:"required,email"`
	Roles           []string `json:"roles" binding:"required"`
	Password        string   `json:"password" binding:"required"`
	PasswordConfirm string   `json:"password_confirm" binding:"eqfield=Password"`
}

func toBusNewUser(app newUser) (userbus.NewUser, error) {
	roles, err := userbus.ParseRoles(app.Roles)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	name, err := userbus.ParseName(app.Name)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := userbus.NewUser{
		Name:     name,
		Email:    *addr,
		Roles:    roles,
		Password: app.Password,
	}

	return bus, nil
}

// =============================================================================

// updateUser defines the data needed to update a user.
type updateUser struct {
	Name            *string `json:"name"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"password_confirm" binding:"omitempty,eqfield=Password"`
}

func toBusUpdateUser(app updateUser) (userbus.UpdateUser, error) {
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
