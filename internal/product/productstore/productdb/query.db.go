package productdb

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

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

	var dbPrds []product
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbPrds); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusProducts(dbPrds)
}
