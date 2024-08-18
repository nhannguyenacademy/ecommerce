package productdb

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

func (s *Store) Create(ctx context.Context, prd productbus.Product) error {
	const q = `
	INSERT INTO products
		(product_id, name, description, image_url, price, quantity, date_created, date_updated)
	VALUES
		(:product_id, :name, :description, :image_url, :price, :quantity, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBProduct(prd)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
