package models

type PaymentMethod struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"type:varchar(255);not null" json:"name"`
	Type           string `gorm:"type:varchar(255);not null" json:"type"`             // e.g., cash, bank_transfer, credit_card
	PaymentMethod  string `gorm:"type:varchar(255)" json:"payment_method,omitempty"`  // e.g., va, qris
	PaymentChannel string `gorm:"type:varchar(255)" json:"payment_channel,omitempty"` // e.g., mandiri, qris
	IsActive       bool   `gorm:"default:true" json:"is_active"`
}
