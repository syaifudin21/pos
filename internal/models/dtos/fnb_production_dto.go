package dtos

import "github.com/google/uuid"

type FNBProductionRequest struct {
	FNBMainProductUuid uuid.UUID `json:"fnb_main_product_uuid" validate:"required"`
	QuantityToProduce  float64   `json:"quantity_to_produce" validate:"required,gt=0"`
}

type FNBProductionResponse struct {
	FNBMainProductUuid uuid.UUID `json:"fnb_main_product_uuid"`
	ProductName        string    `json:"product_name"`
	QuantityProduced   float64   `json:"quantity_produced"`
	Message            string    `json:"message"`
}
