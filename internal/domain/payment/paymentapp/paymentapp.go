package paymentapp

import "github.com/nhannguyenacademy/ecommerce/internal/domain/order/orderbus"

type PaymentPartner interface {
	CreateOrder(amount int64) (string, error)
	QueryOrder(orderID string) (orderbus.Status, error)
}
