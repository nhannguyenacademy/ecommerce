package userdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

// QueryByEmailConfirmToken gets the specified user from the database by email confirm token.
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

	var dbUsr user
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return userbus.User{}, fmt.Errorf("db: %w", userbus.ErrNotFound)
		}
		return userbus.User{}, fmt.Errorf("db: %w", err)
	}

	return toBusUser(dbUsr)
}
