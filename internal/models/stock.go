package models

type Stock struct {
	BaseModel
	OutletID  uint    `gorm:"not null" json:"outlet_id"`
	Outlet    Outlet  `json:"outlet"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `json:"product"`
	Quantity  float64 `gorm:"not null" json:"quantity"`
	UserID    uint `gorm:"not null" json:"user_id"`
	User      User    `json:"user"`
}
