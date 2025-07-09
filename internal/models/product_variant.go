package models

type ProductVariant struct {
	BaseModel
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Product   Product `json:"product"`
	Name      string  `gorm:"not null" json:"name"` // e.g., "Red / L"
	SKU       string  `gorm:"uniqueIndex:idx_user_variant_sku;not null" json:"sku"`
	Price     float64 `gorm:"not null" json:"price"`
	UserID    uint    `gorm:"uniqueIndex:idx_user_variant_sku;not null" json:"user_id"`
}
