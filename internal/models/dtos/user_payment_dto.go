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

type PaymentMethodWithUserStatusResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	PaymentMethod  string `json:"payment_method"`
	PaymentChannel string `json:"payment_channel"`
	Issuer         string `json:"issuer"`
	IsUserActive   bool   `json:"is_user_active"` // Indicates if this payment method is active for the current user
}