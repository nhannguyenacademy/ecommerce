package productdb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

func (s *Store) QueryByIDs(ctx context.Context, prdIDs []uuid.UUID) ([]productbus.Product, error) {
	ids := make([]string, len(prdIDs))
	for i, id := range prdIDs {
		ids[i] = id.String()
	}

	const q = `
	SELECT
        product_id, name, description, image_url, price, quantity, date_created, date_updated
	FROM
		products
	WHERE 
		product_id IN (?)`

	var dbPrds []product
	if err := sqldb.NamedQuerySliceUsingIn(ctx, s.log, s.db, q, ids, &dbPrds); err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	return toBusProducts(dbPrds)
}
