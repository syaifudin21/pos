package models

type OrderItem struct {
	BaseModel
	OrderID   uint    `gorm:"not null" json:"order_id"`
	Order     Order   `json:"order"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `json:"product"`
	Quantity  float64 `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"` // Price at the time of order
}
