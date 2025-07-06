package models

type PaymentMethod struct {
	BaseModel
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Type      string `gorm:"type:varchar(50);not null" json:"type"` // e.g., cash, e-wallet, bank transfer
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	CreatorID *uint  `json:"creator_id"`
	Creator   *User  `gorm:"foreignKey:CreatorID" json:"creator"`
}
