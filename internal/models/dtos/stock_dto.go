package dtos

import "github.com/google/uuid"

type UpdateStockRequest struct {
	ProductUuid        uuid.UUID `json:"product_uuid,omitempty"`
	ProductVariantUuid uuid.UUID `json:"product_variant_uuid,omitempty"`
	Quantity           float64   `json:"quantity" validate:"required"`
}

type GlobalStockUpdateRequest struct {
	OutletUuid         uuid.UUID `json:"outlet_uuid" validate:"required"`
	ProductUuid        uuid.UUID `json:"product_uuid,omitempty"`
	ProductVariantUuid uuid.UUID `json:"product_variant_uuid,omitempty"`
	Quantity           float64   `json:"quantity" validate:"required"`
}

type StockResponse struct {
	ProductUuid        uuid.UUID  `json:"product_uuid"`
	ProductName        string     `json:"product_name"`
	ProductSku         string     `json:"product_sku,omitempty"`
	ProductType        string     `json:"product_type,omitempty"`
	ProductVariantUuid *uuid.UUID `json:"product_variant_uuid,omitempty"`
	VariantName        string     `json:"variant_name,omitempty"`
	VariantSku         string     `json:"variant_sku,omitempty"`
	Quantity           float64    `json:"quantity"`
}

type StockDetailResponse struct {
	Uuid        uuid.UUID                `json:"uuid"`
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Price       float64                  `json:"price"`
	SKU         string                   `json:"sku,omitempty"`
	Type        string                   `json:"type"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
	Recipes     []RecipeResponse         `json:"recipes,omitempty"`
	AddOns      []ProductAddOnResponse   `json:"add_ons,omitempty"`
	Quantity    float64                  `json:"quantity"`
}
