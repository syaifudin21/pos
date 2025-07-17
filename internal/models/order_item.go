package models

type OrderItem struct {
	BaseModel
	OrderID          uint            `gorm:"not null" json:"order_id"`
	Order            Order           `json:"order"`
	ProductID        *uint           `gorm:"index" json:"product_id,omitempty"`
	Product          *Product        `json:"product,omitempty"`
	ProductVariantID *uint           `gorm:"index" json:"product_variant_id,omitempty"`
	ProductVariant   *ProductVariant `json:"product_variant,omitempty"`
	Quantity         float64         `gorm:"not null" json:"quantity"`
	Price            float64         `gorm:"not null" json:"price"` // Price at the time of order
	ProductName      string          `gorm:"type:varchar(255)" json:"product_name"`
	AddOns           []OrderItemAddOn `gorm:"foreignKey:OrderItemID" json:"add_ons,omitempty"`
}
