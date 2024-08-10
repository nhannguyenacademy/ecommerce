package userdb

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type user struct {
	ID                uuid.UUID                `db:"user_id"`
	Name              string                   `db:"name"`
	Email             string                   `db:"email"`
	Roles             pgtype.FlatArray[string] `db:"roles"`
	PasswordHash      string                   `db:"password_hash"`
	Enabled           bool                     `db:"enabled"`
	EmailConfirmToken pgtype.Text              `db:"email_confirm_token"`
	DateCreated       time.Time                `db:"date_created"`
	DateUpdated       time.Time                `db:"date_updated"`
}

func toDBUser(bus userbus.User) user {
	return user{
		ID:                bus.ID,
		Name:              bus.Name.String(),
		Email:             bus.Email.Address,
		Roles:             userbus.ParseRolesToString(bus.Roles),
		PasswordHash:      bus.PasswordHash,
		Enabled:           bus.Enabled,
		EmailConfirmToken: pgtype.Text{String: bus.EmailConfirmToken, Valid: bus.EmailConfirmToken != ""},
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
		DateCreated:       db.DateCreated.In(time.Local),
		DateUpdated:       db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusUsers(dbs []user) ([]userbus.User, error) {
	bus := make([]userbus.User, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusUser(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
