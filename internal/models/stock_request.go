package models

import "github.com/google/uuid"

type GlobalStockUpdateRequest struct {
	OutletUuid  uuid.UUID `json:"outlet_uuid"`
	Productuuid uuid.UUID `json:"product_uuid"`
	Quantity    float64   `json:"quantity"`
}
