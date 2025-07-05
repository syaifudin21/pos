package models

type Recipe struct {
	BaseModel
	MainProductID uint    `gorm:"not null" json:"main_product_id"`
	MainProduct   Product `gorm:"foreignKey:MainProductID" json:"main_product"`
	ComponentID   uint    `gorm:"not null" json:"component_id"`
	Component     Product `gorm:"foreignKey:ComponentID" json:"component"`
	Quantity      float64 `gorm:"not null" json:"quantity"` // Quantity of component needed for one main product
	UserID        uint    `gorm:"not null" json:"user_id"`
	User          User    `json:"user"`
}
