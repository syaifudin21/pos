package dtos

import "github.com/google/uuid"

type UpdateStockRequest struct {
	Quantity float64 `json:"quantity" validate:"required"`
}

type GlobalStockUpdateRequest struct {
	OutletUuid  uuid.UUID `json:"outlet_uuid" validate:"required"`
	Productuuid uuid.UUID `json:"product_uuid" validate:"required"`
	Quantity    float64   `json:"quantity" validate:"required"`
}

type StockResponse struct {
	ProductUuid uuid.UUID `json:"product_uuid"`
	ProductName string    `json:"product_name"`
	ProductSku  string    `json:"product_sku"`
	Quantity    float64   `json:"quantity"`
}
