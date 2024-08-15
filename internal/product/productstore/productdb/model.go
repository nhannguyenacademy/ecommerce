package productdb

import (
	"database/sql"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type product struct {
	ID          uuid.UUID      `db:"product_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	ImageURL    sql.NullString `db:"image_url"`
	Price       int64          `db:"price"`
	Quantity    int32          `db:"quantity"`
	DateCreated time.Time      `db:"date_created"`
	DateUpdated time.Time      `db:"date_updated"`
}

func toDBProduct(bus productbus.Product) product {
	return product{
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

func toBusProduct(db product) (productbus.Product, error) {
	name, err := productbus.ParseName(db.Name)
	if err != nil {
		return productbus.Product{}, fmt.Errorf("parse name: %w", err)
	}

	var imageURL url.URL
	if db.ImageURL.Valid {
		imageURLPtr, err := url.Parse(db.ImageURL.String)
		if err != nil {
			return productbus.Product{}, fmt.Errorf("parse url: %w", err)
		}
		imageURL = *imageURLPtr
	}

	bus := productbus.Product{
		ID:          db.ID,
		Name:        name,
		Description: db.Description.String,
		ImageURL:    imageURL,
		Price:       db.Price,
		Quantity:    db.Quantity,
		DateCreated: db.DateCreated.UTC(),
		DateUpdated: db.DateUpdated.UTC(),
	}

	return bus, nil
}

func toBusProducts(dbs []product) ([]productbus.Product, error) {
	bus := make([]productbus.Product, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusProduct(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
