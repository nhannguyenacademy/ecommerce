package productdb

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

func (s *Store) Update(ctx context.Context, prd productbus.Product) error {
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

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBProduct(prd)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
