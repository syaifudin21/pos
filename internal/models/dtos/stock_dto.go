package dtos

import "github.com/google/uuid"

type UpdateStockRequest struct {
	Quantity float64 `json:"quantity"`
}

type GlobalStockUpdateRequest struct {
	OutletUuid  uuid.UUID `json:"outlet_uuid"`
	Productuuid uuid.UUID `json:"product_uuid"`
	Quantity    float64   `json:"quantity"`
}

type StockResponse struct {
	ProductUuid uuid.UUID `json:"product_uuid"`
	ProductName string    `json:"product_name"`
	ProductSku  string    `json:"product_sku"`
	Quantity    float64   `json:"quantity"`
}
