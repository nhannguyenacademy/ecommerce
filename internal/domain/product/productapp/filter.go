package productapp

import (
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"strconv"
	"time"
)

func parseFilter(qp queryParams) (productbus.QueryFilter, error) {
	var filter productbus.QueryFilter

	if qp.Name != "" {
		name, err := productbus.ParseName(qp.Name)
		if err != nil {
			return productbus.QueryFilter{}, fmt.Errorf("parse name: %w", err)
		}
		filter.Name = &name
	}

	if qp.StartCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartCreatedDate)
		if err != nil {
			return productbus.QueryFilter{}, fmt.Errorf("parse start_created_date: %w", err)
		}
		filter.StartCreatedDate = &t
	}

	if qp.EndCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndCreatedDate)
		if err != nil {
			return productbus.QueryFilter{}, fmt.Errorf("parse end_created_date: %w", err)
		}
		filter.EndCreatedDate = &t
	}

	if qp.StartPrice != "" {
		startPrice, err := strconv.ParseInt(qp.StartPrice, 10, 64)
		if err != nil {
			return productbus.QueryFilter{}, fmt.Errorf("parse start_price: %w", err)
		}
		filter.StartPrice = &startPrice
	}

	if qp.EndPrice != "" {
		endPrice, err := strconv.ParseInt(qp.EndPrice, 10, 64)
		if err != nil {
			return productbus.QueryFilter{}, fmt.Errorf("parse end_price: %w", err)
		}
		filter.EndPrice = &endPrice
	}

	return filter, nil
}
