package orderbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"time"
)

func (b *Business) Create(ctx context.Context, input NewOrder) (Order, error) {
	prodItemMap := make(map[uuid.UUID]NewOrderItem)
	prodIDs := make([]uuid.UUID, len(input.Items))
	for i, itm := range input.Items {
		if _, exist := prodItemMap[itm.ProductID]; exist {
			return Order{}, fmt.Errorf("duplicate product id: %s", itm.ProductID)
		}
		prodItemMap[itm.ProductID] = itm
		prodIDs[i] = itm.ProductID
	}

	prds, err := b.productBus.QueryByIDs(ctx, prodIDs)
	if err != nil {
		return Order{}, fmt.Errorf("query products by ids: %w", err)
	}
	if len(prds) != len(prodIDs) {
		return Order{}, ErrMissingProducts
	}

	var (
		orderAmount int64
		ordID       = uuid.New()
		ordItms     = make([]OrderItem, len(prds))
		now         = time.Now()
	)
	for i, prd := range prds {
		if prodItemMap[prd.ID].Quantity > prd.Quantity {
			return Order{}, fmt.Errorf("insufficient quantity for product: %s", prd.ID)
		}

		remaining := prd.Quantity - prodItemMap[prd.ID].Quantity
		if _, err = b.productBus.Update(ctx, prd.ID, productbus.UpdateProduct{
			Quantity: &remaining,
		}); err != nil {
			return Order{}, fmt.Errorf("update product [%s]: %w", prd.ID, err)
		}

		ordItms[i] = OrderItem{
			ID:          uuid.New(),
			OrderID:     ordID,
			ProductID:   prd.ID,
			Price:       prd.Price,
			Quantity:    prodItemMap[prd.ID].Quantity,
			DateCreated: now,
			DateUpdated: now,
		}

		orderAmount += prd.Price * int64(prodItemMap[prd.ID].Quantity)
	}

	ord := Order{
		ID:          ordID,
		UserID:      input.UserID,
		Amount:      orderAmount,
		Status:      Statuses.Created,
		DateCreated: now,
		DateUpdated: now,
		Items:       ordItms,
	}
	if err := b.storer.Create(ctx, ord); err != nil {
		return Order{}, fmt.Errorf("create: %w", err)
	}

	return ord, nil
}
