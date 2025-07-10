package dtos

import "github.com/google/uuid"

type ProductAddOnRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	AddOnID   uuid.UUID `json:"add_on_id" validate:"required"`
	Price     float64   `json:"price" validate:"required,gt=0"`
}

type ProductAddOnResponse struct {
	ID          uint      `json:"id"`
	Uuid        uuid.UUID `json:"uuid"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	AddOnID     uuid.UUID `json:"add_on_id"`
	AddOnName   string    `json:"add_on_name"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
}
