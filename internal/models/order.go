package models

type Order struct {
	BaseModel
	OutletID      uint        `gorm:"not null" json:"outlet_id"`
	Outlet        Outlet      `json:"outlet"`
	UserID        uint        `gorm:"not null" json:"user_id"`
	User          User        `json:"user"`
	TotalAmount   float64     `gorm:"not null" json:"total_amount"`
	Status        string      `gorm:"not null" json:"status"` // e.g., "pending", "completed", "cancelled"
	PaidAmount    float64     `gorm:"default:0" json:"paid_amount"`
	OrderItems    []OrderItem `json:"order_items"`
	OrderPayments []OrderPayment `json:"order_payments"`
}
