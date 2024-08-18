package userdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"net/mail"
)

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

	var dbUsr user
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return userbus.User{}, fmt.Errorf("db: %w", userbus.ErrNotFound)
		}
		return userbus.User{}, fmt.Errorf("db: %w", err)
	}

	return toBusUser(dbUsr)
}
