package models

type OrderPaymentItem struct {
	BaseModel
	OrderPaymentID uint         `gorm:"not null" json:"order_payment_id"`
	OrderPayment   *OrderPayment `json:"order_payment"`
	OrderItemID    uint         `gorm:"not null" json:"order_item_id"`
	OrderItem      OrderItem    `json:"order_item"`
	QuantityPaid   float64      `gorm:"not null;default:0" json:"quantity_paid"`
}
