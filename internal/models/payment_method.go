package models

type PaymentMethod struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"type:varchar(255);not null" json:"name"`
	Type           string `gorm:"type:varchar(255);not null" json:"type"`             // contoh: cash, bank_transfer, credit_card
	PaymentMethod  string `gorm:"type:varchar(255)" json:"payment_method,omitempty"`  // contoh: va, qris
	PaymentChannel string `gorm:"type:varchar(255)" json:"payment_channel,omitempty"` // contoh: mandiri, qris
	Issuer         string `gorm:"type:varchar(255)" json:"issuer,omitempty"`          // contoh: bca, mandiri, visa, mastercard
	IsActive       bool   `gorm:"default:true" json:"is_active"`
}
