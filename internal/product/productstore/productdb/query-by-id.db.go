package productdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

// QueryByID gets the specified product from the database.
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

	var dbPrd product
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbPrd); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return productbus.Product{}, fmt.Errorf("db: %w", productbus.ErrNotFound)
		}
		return productbus.Product{}, fmt.Errorf("db: %w", err)
	}

	return toBusProduct(dbPrd)
}
