package dtos

type TsmRegisterRequest struct {
	AppCode      string `json:"app_code" validate:"required"`
	MerchantCode string `json:"merchant_code" validate:"required"`
	TerminalCode string `json:"terminal_code" validate:"required"`
	SerialNumber string `json:"serial_number" validate:"required"`
	MID          string `json:"mid" validate:"required"`
	VaIpaymu     string `json:"va_ipaymu"`
}
