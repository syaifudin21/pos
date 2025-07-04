package models

import "github.com/google/uuid"

type CreateOrderRequest struct {
	OutletUuid uuid.UUID          `json:"outlet_uuid"`
	Items      []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductUuid uuid.UUID `json:"product_uuid"`
	Quantity    float64   `json:"quantity"`
}
