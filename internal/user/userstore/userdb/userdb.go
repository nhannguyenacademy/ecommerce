// Package userdb contains user related CRUD functionality.
package userdb

import (
	"github.com/jmoiron/sqlx"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
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
