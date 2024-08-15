package productapp

import (
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"net/http"
	"net/url"
	"time"
)

// =============================================================================
// Query params

// queryParams represents the set of possible query strings.
type queryParams struct {
	Page             string
	Rows             string
	OrderBy          string
	Name             string
	StartCreatedDate string
	EndCreatedDate   string
	StartPrice       string
	EndPrice         string
}

func parseQueryParams(r *http.Request) queryParams {
	values := r.URL.Query()

	filter := queryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("row"),
		OrderBy:          values.Get("order_by"),
		Name:             values.Get("name"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
		StartPrice:       values.Get("start_price"),
		EndPrice:         values.Get("end_price"),
	}

	return filter
}

// =============================================================================

// product represents information about an individual product.
type product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Price       int64  `json:"price"`
	Quantity    int32  `json:"quantity"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

func toAppProduct(bus productbus.Product) product {
	return product{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		Description: bus.Description,
		ImageURL:    bus.ImageURL.String(),
		Price:       bus.Price,
		Quantity:    bus.Quantity,
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProducts(prds []productbus.Product) []product {
	app := make([]product, len(prds))
	for i, usr := range prds {
		app[i] = toAppProduct(usr)
	}

	return app
}

// =============================================================================

type newProduct struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	Price       int64  `json:"price" binding:"required,gte=1"`
	Quantity    int32  `json:"quantity" binding:"required,gte=1"`
}

func toBusNewProduct(app newProduct) (productbus.NewProduct, error) {
	imageURL, err := url.Parse(app.ImageURL)
	if err != nil {
		return productbus.NewProduct{}, fmt.Errorf("parse: %w", err)
	}

	name, err := productbus.ParseName(app.Name)
	if err != nil {
		return productbus.NewProduct{}, fmt.Errorf("parse: %w", err)
	}

	bus := productbus.NewProduct{
		Name:        name,
		Description: app.Description,
		ImageURL:    *imageURL,
		Price:       app.Price,
		Quantity:    app.Quantity,
	}

	return bus, nil
}

// =============================================================================

type updateProduct struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url" binding:"omitempty,url"`
	Price       *int64  `json:"price" binding:"omitempty,gte=1"`
	Quantity    *int32  `json:"quantity" binding:"omitempty,gte=1"`
}

func toBusUpdateProduct(app updateProduct) (productbus.UpdateProduct, error) {
	var name *productbus.Name
	if app.Name != nil {
		nm, err := productbus.ParseName(*app.Name)
		if err != nil {
			return productbus.UpdateProduct{}, fmt.Errorf("parse: %w", err)
		}
		name = &nm
	}

	var imageURL *url.URL
	if app.ImageURL != nil {
		imgURL, err := url.Parse(*app.ImageURL)
		if err != nil {
			return productbus.UpdateProduct{}, fmt.Errorf("parse: %w", err)
		}
		imageURL = imgURL
	}

	bus := productbus.UpdateProduct{
		Name:        name,
		Description: app.Description,
		ImageURL:    imageURL,
		Price:       app.Price,
		Quantity:    app.Quantity,
	}

	return bus, nil
}
