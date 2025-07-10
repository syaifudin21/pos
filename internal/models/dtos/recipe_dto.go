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
	Uuid          uuid.UUID `json:"uuid"`
	ComponentName string    `json:"component_name"`
	Quantity      float64   `json:"quantity"`
}
