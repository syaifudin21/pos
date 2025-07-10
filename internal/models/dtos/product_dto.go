package dtos

import "github.com/google/uuid"

type ProductVariantCreateRequest struct {
	Name  string  `json:"name" validate:"required"`
	SKU   string  `json:"sku" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

type ProductVariantUpdateRequest struct {
	ID    uint    `json:"id,omitempty"` // Include ID for updating existing variants
	Name  string  `json:"name" validate:"required"`
	SKU   string  `json:"sku" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

type ProductVariantResponse struct {
	ID    uint      `json:"id"`
	Uuid  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	SKU   string    `json:"sku"`
	Price float64   `json:"price"`
}

type ProductCreateRequest struct {
	Name        string                        `json:"name" validate:"required"`
	Description string                        `json:"description,omitempty"`
	Price       float64                       `json:"price"`
	SKU         string                        `json:"sku,omitempty"`
	Type        string                        `json:"type" validate:"required,oneof=retail_item fnb_main_product fnb_component add_on"`
	Variants    []ProductVariantCreateRequest `json:"variants,omitempty"`
}

type ProductUpdateRequest struct {
	Name        string                        `json:"name" validate:"required"`
	Description string                        `json:"description,omitempty"`
	Price       float64                       `json:"price"`
	SKU         string                        `json:"sku,omitempty"`
	Type        string                        `json:"type" validate:"required,oneof=retail_item fnb_main_product fnb_component add_on"`
	Variants    []ProductVariantUpdateRequest `json:"variants,omitempty"`
}

type ProductResponse struct {
	ID          uint                     `json:"id"`
	Uuid        uuid.UUID                `json:"uuid"`
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Price       float64                  `json:"price"`
	SKU         string                   `json:"sku,omitempty"`
	Type        string                   `json:"type"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
	Recipes     []RecipeResponse         `json:"recipes,omitempty"`
	AddOns      []ProductAddOnResponse   `json:"add_ons,omitempty"`
}

type ProductOutletResponse struct {
	ProductUuid uuid.UUID `json:"product_uuid"`
	ProductName string    `json:"product_name"`
	ProductSku  string    `json:"product_sku"`
	Price       float64   `json:"price"`
	Type        string    `json:"type"`
	Quantity    float64   `json:"quantity"` // Stock quantity at the outlet
}
