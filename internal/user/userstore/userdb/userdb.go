// Package userdb contains user related CRUD functionality.
package userdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"net/mail"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

func (s *Store) Create(ctx context.Context, user userbus.User) error {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, enabled, email_confirm_token, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :enabled, :email_confirm_token, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBUser(user)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", userbus.ErrUniqueEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, user userbus.User) error {
	const q = `
	UPDATE
		users
	SET 
		"name" = :name,
		"email" = :email,
		"roles" = :roles,
		"password_hash" = :password_hash,
		"enabled" = :enabled,
		"email_confirm_token" = :email_confirm_token,
		"date_updated" = :date_updated
	WHERE
		user_id = :user_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBUser(user)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return userbus.ErrUniqueEmail
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (userbus.User, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
        user_id, name, email, password_hash, roles, enabled, email_confirm_token, date_created, date_updated
	FROM
		users
	WHERE 
		user_id = :user_id`

	var row userRow
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return userbus.User{}, fmt.Errorf("db: %w", userbus.ErrNotFound)
		}
		return userbus.User{}, fmt.Errorf("db: %w", err)
	}

	return toBusUser(row)
}

func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (userbus.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email.Address,
	}

	const q = `
	SELECT
        user_id, name, email, password_hash, roles, enabled, email_confirm_token, date_created, date_updated
	FROM
		users
	WHERE
		email = :email`

	var row userRow
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return userbus.User{}, fmt.Errorf("db: %w", userbus.ErrNotFound)
		}
		return userbus.User{}, fmt.Errorf("db: %w", err)
	}

	return toBusUser(row)
}

func (s *Store) QueryByEmailConfirmToken(ctx context.Context, emailConfirmToken string) (userbus.User, error) {
	data := struct {
		EmailConfirmToken string `db:"email_confirm_token"`
	}{
		EmailConfirmToken: emailConfirmToken,
	}

	const q = `
	SELECT
        user_id, name, email, password_hash, roles, enabled, email_confirm_token, date_created, date_updated
	FROM
		users
	WHERE
		email_confirm_token = :email_confirm_token`

	var row userRow
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return userbus.User{}, fmt.Errorf("db: %w", userbus.ErrNotFound)
		}
		return userbus.User{}, fmt.Errorf("db: %w", err)
	}

	return toBusUser(row)
}
