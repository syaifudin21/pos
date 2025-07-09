package models

// AllowedProductTypes defines the list of types that a product can have.
var AllowedProductTypes = []string{"retail_item", "fnb_main_product", "fnb_component"}

type Product struct {
	BaseModel
	Name        string           `gorm:"not null" json:"name"`
	Description string           `json:"description,omitempty"`
	Price       float64          `gorm:"not null" json:"price"`
	SKU         string           `gorm:"uniqueIndex:idx_user_sku" json:"sku,omitempty"`
	Type        string           `gorm:"not null" json:"type"` // e.g., "retail_item", "fnb_main_product", "fnb_component"
	UserID      uint             `gorm:"uniqueIndex:idx_user_sku;not null" json:"user_id"`
	User        User             `json:"user"`
	Variants    []ProductVariant `json:"variants,omitempty"`
	Recipes     []Recipe         `gorm:"foreignKey:MainProductID" json:"recipes,omitempty"`
}
