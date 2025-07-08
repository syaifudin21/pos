package dtos

type ActivateUserPaymentRequest struct {
	PaymentMethodID uint `json:"payment_method_id" validate:"required"`
}

type DeactivateUserPaymentRequest struct {
	PaymentMethodID uint `json:"payment_method_id" validate:"required"`
}

type UserPaymentResponse struct {
	PaymentMethodID uint   `json:"payment_method_id"`
	PaymentName     string `json:"payment_name"`
	PaymentMethod   string `json:"payment_method"`
	PaymentChannel  string `json:"payment_channel"`
	IsActive        bool   `json:"is_active"`
}
