package commands

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/user/userstore/userdb"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"net/mail"
	"time"
)

// UserAdd adds new users into the database.
func UserAdd(log *logger.Logger, cfg sqldb.Config, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userBus := userbus.NewBusiness(log, userdb.NewStore(log, db))

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("parsing email: %w", err)
	}

	nu := userbus.NewUser{
		Name:     userbus.MustParseName(name),
		Email:    *addr,
		Password: password,
		Roles:    []userbus.Role{userbus.Roles.Admin, userbus.Roles.User},
	}

	usr, err := userBus.Create(ctx, nu)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
