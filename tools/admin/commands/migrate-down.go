package commands

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/migrate"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"time"
)

func MigrateDown(cfg sqldb.Config) error {
	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrate.MigrateDown(ctx, db); err != nil {
		return fmt.Errorf("migrate down database: %w", err)
	}

	fmt.Println("migrate down complete")
	return nil
}
