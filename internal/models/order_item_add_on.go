package models

type OrderItemAddOn struct {
	BaseModel
	OrderItemID uint    `gorm:"not null;index" json:"order_item_id"`
	OrderItem   OrderItem `json:"order_item"`
	AddOnID     uint    `gorm:"not null;index" json:"add_on_id"`
	AddOn       Product `json:"add_on"` // The add-on product itself
	Quantity    float64 `gorm:"not null" json:"quantity"`
	Price       float64 `gorm:"not null" json:"price"` // Price of the add-on at the time of order
	UserID      uint    `gorm:"not null" json:"user_id"`
	User        User    `json:"user"`
}
