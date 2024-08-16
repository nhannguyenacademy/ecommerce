package productdb

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

// Delete removes a product from the database.
func (s *Store) Delete(ctx context.Context, prd productbus.Product) error {
	const q = `
	DELETE FROM
		products
	WHERE
		product_id = :product_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBProduct(prd)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
