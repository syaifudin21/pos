package models

import "github.com/google/uuid"

type GlobalStockUpdateRequest struct {
	OutletExternalID  uuid.UUID `json:"outlet_uuid"`
	ProductExternalID uuid.UUID `json:"product_uuid"`
	Quantity          float64   `json:"quantity"`
}
