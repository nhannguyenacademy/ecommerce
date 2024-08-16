package productdb

import (
	"bytes"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"strings"
)

func applyFilter(filter productbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.StartCreatedDate != nil {
		data["start_date_created"] = filter.StartCreatedDate.UTC()
		wc = append(wc, "date_created >= :start_date_created")
	}

	if filter.EndCreatedDate != nil {
		data["end_date_created"] = filter.EndCreatedDate.UTC()
		wc = append(wc, "date_created <= :end_date_created")
	}

	if filter.StartPrice != nil {
		data["start_price"] = filter.StartPrice
		wc = append(wc, "price >= :start_price")
	}

	if filter.EndPrice != nil {
		data["end_price"] = filter.EndPrice
		wc = append(wc, "price <= :end_price")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
