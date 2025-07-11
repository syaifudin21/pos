package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderPaymentRequest struct {
	OrderUuid       uuid.UUID `json:"order_uuid" validate:"required"`
	PaymentMethodID uint      `json:"payment_method_id" validate:"required"`
	AmountPaid      float64   `json:"amount_paid" validate:"required,gt=0"`
	CustomerName    string    `json:"customer_name"`
	CustomerPhone   string    `json:"customer_phone"`
}

type OrderPaymentResponse struct {
	Uuid            uuid.UUID  `json:"uuid"`
	OrderUuid       uuid.UUID  `json:"order_uuid"`
	PaymentMethodID uint       `json:"payment_method_id"`
	PaymentName     string     `json:"payment_name"`
	AmountPaid      float64    `json:"amount_paid"`
	CustomerName    string     `json:"customer_name"`
	CustomerPhone   string     `json:"customer_phone"`
	ChangeAmount    float64    `json:"change_amount"`
	CreatedAt       string     `json:"created_at"`
	IsPaid          bool       `json:"is_paid"` // This might be derived or from a new field in OrderPayment model
	PaidAt          *time.Time `json:"paid_at"` // Use pointer for nullable timestamp

}
