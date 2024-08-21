// Package productdb contains product related CRUD functionality.
package productdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

// Store manages the set of APIs for database access.
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

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (productbus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

func (s *Store) Create(ctx context.Context, product productbus.Product) error {
	const q = `
	INSERT INTO products
		(product_id, name, description, image_url, price, quantity, date_created, date_updated)
	VALUES
		(:product_id, :name, :description, :image_url, :price, :quantity, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBProduct(product)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, product productbus.Product) error {
	const q = `
	UPDATE
		products
	SET 
		"name" = :name,
		"description" = :description,
		"image_url" = :image_url,
		"price" = :price,
		"quantity" = :quantity,
		"date_updated" = :date_updated
	WHERE
		product_id = :product_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBProduct(product)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByIDs(ctx context.Context, productIDs []uuid.UUID) ([]productbus.Product, error) {
	const q = `
	SELECT
        product_id, name, description, image_url, price, quantity, date_created, date_updated
	FROM
		products
	WHERE 
		product_id IN (:product_ids)`

	inData := map[string]any{
		"product_ids": productIDs,
	}
	var rows []productRow
	if err := sqldb.NamedQuerySliceUsingIn(ctx, s.log, s.db, q, inData, &rows); err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	return toBusProducts(rows)
}

func (s *Store) Delete(ctx context.Context, product productbus.Product) error {
	const q = `
	DELETE FROM
		products
	WHERE
		product_id = :product_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBProduct(product)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Query(ctx context.Context, filter productbus.QueryFilter, sortBy sort.By, page page.Page) ([]productbus.Product, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		product_id, name, description, image_url, price, quantity, date_created, date_updated
	FROM
		products`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(sortBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var rows []productRow
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &rows); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusProducts(rows)
}

func (s *Store) QueryByID(ctx context.Context, prdID uuid.UUID) (productbus.Product, error) {
	data := struct {
		ID string `db:"product_id"`
	}{
		ID: prdID.String(),
	}

	const q = `
	SELECT
        product_id, name, description, image_url, price, quantity, date_created, date_updated
	FROM
		products
	WHERE 
		product_id = :product_id`

	var row productRow
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return productbus.Product{}, fmt.Errorf("db: %w", productbus.ErrNotFound)
		}
		return productbus.Product{}, fmt.Errorf("db: %w", err)
	}

	return toBusProduct(row)
}

func (s *Store) Count(ctx context.Context, filter productbus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		products`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}
