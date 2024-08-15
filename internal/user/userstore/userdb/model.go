package userdb

import (
	"database/sql"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb/dbarray"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type user struct {
	ID                uuid.UUID      `db:"user_id"`
	Name              string         `db:"name"`
	Email             string         `db:"email"`
	Roles             dbarray.String `db:"roles"`
	PasswordHash      string         `db:"password_hash"`
	Enabled           bool           `db:"enabled"`
	EmailConfirmToken sql.NullString `db:"email_confirm_token"`
	DateCreated       time.Time      `db:"date_created"`
	DateUpdated       time.Time      `db:"date_updated"`
}

func toDBUser(bus userbus.User) user {
	return user{
		ID:                bus.ID,
		Name:              bus.Name.String(),
		Email:             bus.Email.Address,
		Roles:             userbus.ParseRolesToString(bus.Roles),
		PasswordHash:      bus.PasswordHash,
		Enabled:           bus.Enabled,
		EmailConfirmToken: sql.NullString{String: bus.EmailConfirmToken, Valid: bus.EmailConfirmToken != ""},
		DateCreated:       bus.DateCreated.UTC(),
		DateUpdated:       bus.DateUpdated.UTC(),
	}
}

func toBusUser(db user) (userbus.User, error) {
	addr := mail.Address{
		Address: db.Email,
	}

	roles, err := userbus.ParseRoles(db.Roles)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse: %w", err)
	}

	name, err := userbus.ParseName(db.Name)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse name: %w", err)
	}

	bus := userbus.User{
		ID:                db.ID,
		Name:              name,
		Email:             addr,
		Roles:             roles,
		PasswordHash:      db.PasswordHash,
		Enabled:           db.Enabled,
		EmailConfirmToken: db.EmailConfirmToken.String,
		DateCreated:       db.DateCreated.UTC(),
		DateUpdated:       db.DateUpdated.UTC(),
	}

	return bus, nil
}
