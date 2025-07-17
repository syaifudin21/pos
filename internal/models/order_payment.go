package models

import "time"

type OrderPayment struct {
	BaseModel
	OrderID         uint          `gorm:"not null" json:"order_id"`
	Order           Order         `json:"order"`
	PaymentMethodID uint          `gorm:"not null" json:"payment_method_id"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	AmountPaid      float64       `gorm:"not null;column:amount_paid;default:0" json:"amount_paid"`
	ReferenceID     string        `gorm:"type:varchar(255)" json:"reference_id"`
	IsPaid          bool          `gorm:"default:true" json:"is_paid"`
	PaidAt          *time.Time    `json:"paid_at"`
	CustomerName    string        `gorm:"type:varchar(255)" json:"customer_name"`
	CustomerEmail   string        `gorm:"type:varchar(255)" json:"customer_email"`
	CustomerPhone   string        `gorm:"type:varchar(255)" json:"customer_phone"`
	ChangeAmount    float64       `gorm:"default:0" json:"change_amount"`
	Extra           string        `gorm:"type:jsonb" json:"extra,omitempty"`
}
