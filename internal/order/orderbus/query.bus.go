package orderbus

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
)

func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Order, error) {
	ords, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return ords, nil
}
