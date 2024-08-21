package productbus

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
)

func (b *Business) Query(ctx context.Context, filter QueryFilter, sortBy sort.By, page page.Page) ([]Product, error) {
	prds, err := b.storer.Query(ctx, filter, sortBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}
