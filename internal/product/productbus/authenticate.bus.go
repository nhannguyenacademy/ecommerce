package productbus

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

// Authenticate finds a user by their email and verifies their passworb. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (b *Business) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {
	usr, err := b.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, fmt.Errorf("query: email[%s]: %w", email, ErrAuthenticationFailure)
		}
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}
