package dtos

import "github.com/google/uuid"

type CreateRecipeRequest struct {
	MainProductUuid uuid.UUID `json:"main_product_uuid"`
	ComponentUuid   uuid.UUID `json:"component_uuid"`
	Quantity        float64   `json:"quantity"`
}

type UpdateRecipeRequest struct {
	MainProductUuid uuid.UUID `json:"main_product_uuid"`
	ComponentUuid   uuid.UUID `json:"component_uuid"`
	Quantity        float64   `json:"quantity"`
}

type RecipeResponse struct {
	ID              uint      `json:"id"`
	Uuid            uuid.UUID `json:"uuid"`
	MainProductID   uint      `json:"main_product_id"`
	MainProductUuid uuid.UUID `json:"main_product_uuid"`
	ComponentID     uint      `json:"component_id"`
	ComponentUuid   uuid.UUID `json:"component_uuid"`
	Quantity        float64   `json:"quantity"`
}
