package models

type Stock struct {
	BaseModel
	OutletID         uint            `gorm:"not null;index" json:"outlet_id"`
	Outlet           Outlet          `json:"outlet"`
	ProductID        *uint           `gorm:"index" json:"product_id,omitempty"`
	Product          *Product        `json:"product,omitempty"`
	ProductVariantID *uint           `gorm:"index" json:"product_variant_id,omitempty"`
	ProductVariant   *ProductVariant `json:"product_variant,omitempty"`
	Quantity         float64         `gorm:"not null" json:"quantity"`
	UserID           uint            `gorm:"not null" json:"user_id"`
	User             User            `json:"user"`
}
