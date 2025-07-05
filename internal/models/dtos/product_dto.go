package dtos

import "github.com/google/uuid"

type ProductCreateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description,omitempty" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	SKU         string  `json:"sku,omitempty" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=retail_item fnb_main_product fnb_component"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description,omitempty" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	SKU         string  `json:"sku,omitempty" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=retail_item fnb_main_product fnb_component"`
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku,omitempty"`
	Type        string    `json:"type"`
}

type ProductOutletResponse struct {
	ProductUuid uuid.UUID `json:"product_uuid"`
	ProductName string    `json:"product_name"`
	ProductSku  string    `json:"product_sku"`
	Price       float64   `json:"price"`
	Type        string    `json:"type"`
	Quantity    float64   `json:"quantity"` // Stock quantity at the outlet
}
