package models

type ProductAddOn struct {
	BaseModel
	ProductID   uint    `gorm:"not null;index" json:"product_id"` // The main product this add-on belongs to
	Product     Product `json:"product"`
	AddOnID     uint    `gorm:"not null;index" json:"add_on_id"`     // The actual add-on product (type: add_on)
	AddOn       Product `json:"add_on"`
	Price       float64 `gorm:"not null" json:"price"`
	IsAvailable bool    `gorm:"default:true" json:"is_available"`
	UserID      uint    `gorm:"not null" json:"user_id"`
	User        User    `json:"user"`
}
