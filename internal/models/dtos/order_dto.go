package dtos

import "github.com/google/uuid"

type CreateOrderRequest struct {
	OutletUuid    uuid.UUID          `json:"outlet_uuid" validate:"required"`
	Items         []OrderItemRequest `json:"items" validate:"required,dive"`
	PaymentMethod string             `json:"payment_method" validate:"required"`
}

type OrderItemRequest struct {
	ProductUuid        uuid.UUID            `json:"product_uuid,omitempty"`
	ProductVariantUuid uuid.UUID            `json:"product_variant_uuid,omitempty"`
	Quantity           int                  `json:"quantity" validate:"required,gt=0"`
	AddOns             []OrderItemAddonRequest `json:"add_ons,omitempty"`
}

type OrderItemAddonRequest struct {
	AddOnUuid uuid.UUID `json:"add_on_uuid" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}

// OrderResponse represents the response structure for an order.
// This can be a simplified version of models.Order if not all fields are needed.
type OrderResponse struct {
	ID            uint      `json:"id"`
	Uuid          uuid.UUID `json:"uuid"`
	OutletID      uint      `json:"outlet_id"`
	OutletUuid    uuid.UUID `json:"outlet_uuid"`
	UserID        uint      `json:"user_id"`
	UserUuid      uuid.UUID `json:"user_uuid"`
	OrderDate     string    `json:"order_date"` // Consider formatting time.Time to string
	TotalAmount   float64   `json:"total_amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	// Add other fields if necessary, but keep it minimal
}

type OrderItemResponse struct {
	ID          uint      `json:"id"`
	Uuid        uuid.UUID `json:"uuid"`
	OrderID     uint      `json:"order_id"`
	OrderUuid   uuid.UUID `json:"order_uuid"`
	ProductID   uint      `json:"product_id"`
	ProductUuid uuid.UUID `json:"product_uuid"`
	ProductName string    `json:"product_name"`
	Quantity    float64   `json:"quantity"`
	Price       float64   `json:"price"`
}
