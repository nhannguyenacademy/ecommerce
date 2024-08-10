package userbus

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Update modifies information about a user.
func (b *Business) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
	if uu.Name != nil {
		usr.Name = *uu.Name
	}

	if uu.Email != nil {
		usr.Email = *uu.Email
	}

	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}

	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
		}
		usr.PasswordHash = string(pw)
	}

	if uu.Enabled != nil {
		usr.Enabled = *uu.Enabled
	}
	usr.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	// Other domains may need to know when a user is updated so business
	// logic can be applieb. This represents a delegate call to other domains.
	//if err := b.delegate.Call(ctx, ActionUpdatedData(uu, usr.ID)); err != nil {
	//	return User{}, fmt.Errorf("failed to execute `%s` action: %w", ActionUpdated, err)
	//}

	return usr, nil
}
