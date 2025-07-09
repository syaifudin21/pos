package dtos

import "github.com/google/uuid"

type CreateRecipeRequest struct {
	MainProductUuid uuid.UUID `json:"main_product_uuid" validate:"required"`
	ComponentUuid   uuid.UUID `json:"component_uuid" validate:"required"`
	Quantity        float64   `json:"quantity" validate:"required"`
}

type UpdateRecipeRequest struct {
	MainProductUuid uuid.UUID `json:"main_product_uuid" validate:"required"`
	ComponentUuid   uuid.UUID `json:"component_uuid" validate:"required"`
	Quantity        float64   `json:"quantity" validate:"required"`
}

type RecipeResponse struct {
	ID              uint      `json:"id"`
	Uuid            uuid.UUID `json:"uuid"`
	MainProductID   uint      `json:"main_product_id"`
	MainProductUuid uuid.UUID `json:"main_product_uuid"`
	MainProductName string    `json:"main_product_name"` // Added
	ComponentID     uint      `json:"component_id"`
	ComponentUuid   uuid.UUID `json:"component_uuid"`
	ComponentName   string    `json:"component_name"`   // Added
	Quantity        float64   `json:"quantity"`
}
