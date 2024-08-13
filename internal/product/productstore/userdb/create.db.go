package userdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr userbus.User) error {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, enabled, email_confirm_token, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :enabled, :email_confirm_token, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", userbus.ErrUniqueEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
