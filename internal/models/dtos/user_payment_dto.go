package dtos

type ActivateUserPaymentRequest struct {
	PaymentMethodID uint `json:"payment_method_id" validate:"required"`
}

type DeactivateUserPaymentRequest struct {
	PaymentMethodID uint `json:"payment_method_id" validate:"required"`
}