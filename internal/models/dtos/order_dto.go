package dtos

import "github.com/google/uuid"
import "time"

type CreateOrderRequest struct {
	OutletUuid      uuid.UUID          `json:"outlet_uuid" validate:"required"`
	Items           []OrderItemRequest `json:"items" validate:"required,dive"`
	PaymentMethodID uint               `json:"payment_method_id" validate:"required"`
}

type OrderItemRequest struct {
	ProductUuid        uuid.UUID               `json:"product_uuid,omitempty"`
	ProductVariantUuid uuid.UUID               `json:"product_variant_uuid,omitempty"`
	Quantity           int                     `json:"quantity" validate:"required,gt=0"`
	AddOns             []OrderItemAddonRequest `json:"add_ons,omitempty"`
}

type OrderItemAddonRequest struct {
	AddOnUuid uuid.UUID `json:"add_on_uuid" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}

// UserDetailResponse for created_by
type UserDetailResponse struct {
	Uuid uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

// OutletDetailResponse for outlet
type OutletDetailResponse struct {
	Uuid    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Contact string    `json:"contact"`
}

// OrderPaymentDetailResponse for payments
type OrderPaymentDetailResponse struct {
	Uuid            uuid.UUID  `json:"uuid"`
	PaymentMethodID uint       `json:"payment_method_id"`
	PaidAmount      float64    `json:"paid_amount"`
	Name            string     `json:"name"` // Payment method name
	PaymentMethod   string     `json:"payment_method"`
	PaymentChannel  string     `json:"payment_channel"`
	IsPaid          bool       `json:"is_paid"` // This might be derived or from a new field in OrderPayment model
	ReferenceID     string     `json:"reference_id"`
	CreatedAt       string     `json:"created_at"`
	PaidAt          *time.Time `json:"paid_at"` // Use pointer for nullable timestamp
	ChangeAmount    float64    `json:"change_amount"`
}

// OrderItemAddonDetailResponse for add_ons within order items
type OrderItemAddonDetailResponse struct {
	Uuid     uuid.UUID `json:"add_on_uuid"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
}

// OrderItemDetailResponse for items
type OrderItemDetailResponse struct {
	ProductUuid        uuid.UUID                      `json:"product_uuid,omitempty"`
	ProductVariantUuid uuid.UUID                      `json:"product_variant_uuid,omitempty"`
	Name               string                         `json:"name"` // Product name
	Quantity           int                            `json:"quantity"`
	AddOns             []OrderItemAddonDetailResponse `json:"add_ons,omitempty"`
}

// OrderResponse represents the comprehensive response structure for an order.
type OrderResponse struct {
	Uuid          uuid.UUID              `json:"uuid"`
	OrderDate     string                 `json:"order_date"`
	TotalAmount   float64                `json:"total_amount"`
	PaidAmount    float64                `json:"paid_amount"`
	PaymentMethod string                 `json:"payment_method"`
	Status        string                 `json:"status"`
	CreatedBy     *UserDetailResponse     `json:"created_by"`
	Outlet        OutletDetailResponse   `json:"outlet"`
	Payments      []OrderPaymentDetailResponse `json:"payments"`
	Items         []OrderItemDetailResponse    `json:"items"`
}