package paymentbus

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID                   uuid.UUID
	OrderID              uuid.UUID
	Partner              string
	PartnerOrderID       string
	PartnerTransactionID string
	Status               Status
	Currency             string
	DateCreated          time.Time
	DateUpdated          time.Time
}

type NewPayment struct {
	OrderID        uuid.UUID
	Partner        string
	PartnerOrderID string
}

type UpdatePayment struct {
	PartnerTransactionID *string
	Status               *Status
}
