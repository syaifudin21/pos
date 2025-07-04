package models

// AllowedProductTypes defines the list of types that a product can have.
var AllowedProductTypes = []string{"retail_item", "fnb_main_product", "fnb_component"}

type Product struct {
	BaseModel
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `gorm:"not null" json:"price"`
	SKU         string  `gorm:"unique" json:"sku,omitempty"`
	Type        string  `gorm:"not null" json:"type"` // e.g., "retail_item", "fnb_main_product", "fnb_component"
}
