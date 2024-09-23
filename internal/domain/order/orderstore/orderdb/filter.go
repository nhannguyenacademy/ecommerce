package orderdb

import (
	"bytes"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"strings"
)

func applyFilter(filter orderbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.StartCreatedDate != nil {
		data["start_date_created"] = filter.StartCreatedDate.UTC()
		wc = append(wc, "date_created >= :start_date_created")
	}

	if filter.EndCreatedDate != nil {
		data["end_date_created"] = filter.EndCreatedDate.UTC()
		wc = append(wc, "date_created <= :end_date_created")
	}

	if filter.UserID != nil {
		data["user_id"] = filter.UserID
		wc = append(wc, "user_id = :user_id")
	}

	if filter.Status != nil {
		data["status"] = filter.Status.String()
		wc = append(wc, "status = :status")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
