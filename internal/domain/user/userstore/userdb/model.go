package userdb

import (
	"database/sql"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sqldb/dbarray"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type userRow struct {
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

func toDBUser(bus userbus.User) userRow {
	return userRow{
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

func toBusUser(row userRow) (userbus.User, error) {
	addr := mail.Address{
		Address: row.Email,
	}

	roles, err := userbus.ParseRoles(row.Roles)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse: %w", err)
	}

	name, err := userbus.ParseName(row.Name)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse name: %w", err)
	}

	bus := userbus.User{
		ID:                row.ID,
		Name:              name,
		Email:             addr,
		Roles:             roles,
		PasswordHash:      row.PasswordHash,
		Enabled:           row.Enabled,
		EmailConfirmToken: row.EmailConfirmToken.String,
		DateCreated:       row.DateCreated.UTC(),
		DateUpdated:       row.DateUpdated.UTC(),
	}

	return bus, nil
}
