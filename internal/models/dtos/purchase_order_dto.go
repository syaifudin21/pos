package dtos

import "github.com/google/uuid"

type CreatePurchaseOrderRequest struct {
	SupplierUuid uuid.UUID             `json:"supplier_uuid" validate:"required"`
	OutletUuid   uuid.UUID             `json:"outlet_uuid" validate:"required"`
	Items        []PurchaseItemRequest `json:"items" validate:"required,dive"`
}

type PurchaseItemRequest struct {
	ProductUuid        uuid.UUID `json:"product_uuid,omitempty"`
	ProductVariantUuid uuid.UUID `json:"product_variant_uuid,omitempty"`
	Quantity           int       `json:"quantity" validate:"required,gt=0"`
	Price              float64   `json:"price" validate:"required"`
}

type PurchaseOrderResponse struct {
	ID           uint      `json:"id"`
	Uuid         uuid.UUID `json:"uuid"`
	SupplierID   uint      `json:"supplier_id"`
	SupplierUuid uuid.UUID `json:"supplier_uuid"`
	OutletID     uint      `json:"outlet_id"`
	OutletUuid   uuid.UUID `json:"outlet_uuid"`
	OrderDate    string    `json:"order_date"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
}
