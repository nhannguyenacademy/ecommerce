package orderapp

import (
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"time"
)

func parseFilter(qp queryParams) (orderbus.QueryFilter, error) {
	var filter orderbus.QueryFilter

	if qp.StartCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartCreatedDate)
		if err != nil {
			return orderbus.QueryFilter{}, errs.NewFieldsError("start_created_date", err)
		}
		filter.StartCreatedDate = &t
	}

	if qp.EndCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndCreatedDate)
		if err != nil {
			return orderbus.QueryFilter{}, errs.NewFieldsError("end_created_date", err)
		}
		filter.EndCreatedDate = &t
	}

	if qp.UserID != "" {
		id, err := uuid.Parse(qp.UserID)
		if err != nil {
			return orderbus.QueryFilter{}, errs.NewFieldsError("user_id", err)
		}
		filter.UserID = &id
	}

	if qp.Status != "" {
		status, err := orderbus.ParseStatus(qp.Status)
		if err != nil {
			return orderbus.QueryFilter{}, errs.NewFieldsError("status", err)
		}
		filter.Status = &status
	}

	return filter, nil
}
