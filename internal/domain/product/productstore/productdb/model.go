package productdb

import (
	"database/sql"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/product/productbus"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type productRow struct {
	ID          uuid.UUID      `db:"product_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	ImageURL    sql.NullString `db:"image_url"`
	Price       int64          `db:"price"`
	Quantity    int32          `db:"quantity"`
	DateCreated time.Time      `db:"date_created"`
	DateUpdated time.Time      `db:"date_updated"`
}

func toDBProduct(bus productbus.Product) productRow {
	return productRow{
		ID:          bus.ID,
		Name:        bus.Name.String(),
		Description: sql.NullString{String: bus.Description, Valid: bus.Description != ""},
		ImageURL:    sql.NullString{String: bus.ImageURL.String(), Valid: bus.ImageURL.String() != ""},
		Price:       bus.Price,
		Quantity:    bus.Quantity,
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}
}

func toBusProduct(row productRow) (productbus.Product, error) {
	name, err := productbus.ParseName(row.Name)
	if err != nil {
		return productbus.Product{}, fmt.Errorf("parse name: %w", err)
	}

	var imageURL url.URL
	if row.ImageURL.Valid {
		imageURLPtr, err := url.Parse(row.ImageURL.String)
		if err != nil {
			return productbus.Product{}, fmt.Errorf("parse url: %w", err)
		}
		imageURL = *imageURLPtr
	}

	bus := productbus.Product{
		ID:          row.ID,
		Name:        name,
		Description: row.Description.String,
		ImageURL:    imageURL,
		Price:       row.Price,
		Quantity:    row.Quantity,
		DateCreated: row.DateCreated.UTC(),
		DateUpdated: row.DateUpdated.UTC(),
	}

	return bus, nil
}

func toBusProducts(rows []productRow) ([]productbus.Product, error) {
	bus := make([]productbus.Product, len(rows))

	for i, row := range rows {
		var err error
		bus[i], err = toBusProduct(row)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
