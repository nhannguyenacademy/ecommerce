package orderapp

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"time"
)

func parseFilter(qp queryParams) (orderbus.QueryFilter, error) {
	var filter orderbus.QueryFilter

	if qp.StartCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartCreatedDate)
		if err != nil {
			return orderbus.QueryFilter{}, fmt.Errorf("parse start_created_date: %w", err)
		}
		filter.StartCreatedDate = &t
	}

	if qp.EndCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndCreatedDate)
		if err != nil {
			return orderbus.QueryFilter{}, fmt.Errorf("parse end_created_date: %w", err)
		}
		filter.EndCreatedDate = &t
	}

	if qp.UserID != "" {
		id, err := uuid.Parse(qp.UserID)
		if err != nil {
			return orderbus.QueryFilter{}, fmt.Errorf("parse user_id: %w", err)
		}
		filter.UserID = &id
	}

	if qp.Status != "" {
		status, err := orderbus.ParseStatus(qp.Status)
		if err != nil {
			return orderbus.QueryFilter{}, fmt.Errorf("parse status: %w", err)
		}
		filter.Status = &status
	}

	return filter, nil
}
