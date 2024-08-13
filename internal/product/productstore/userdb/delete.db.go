package userdb

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr userbus.User) error {
	const q = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
