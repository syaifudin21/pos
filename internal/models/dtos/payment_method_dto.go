package dtos

type PaymentMethodRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive *bool  `json:"is_active"`
}

type PaymentMethodResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}
