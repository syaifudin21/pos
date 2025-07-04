package models

import "github.com/google/uuid"

type CreatePurchaseOrderRequest struct {
	SupplierUuid uuid.UUID                  `json:"supplier_uuid"`
	OutletUuid   uuid.UUID                  `json:"outlet_uuid"`
	Items        []PurchaseOrderItemRequest `json:"items"`
}

type PurchaseOrderItemRequest struct {
	Productuuid uuid.UUID `json:"product_uuid"`
	Quantity    float64   `json:"quantity"`
	Price       float64   `json:"price"`
}
