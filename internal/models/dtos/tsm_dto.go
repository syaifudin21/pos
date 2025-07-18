package dtos

type TsmRegisterRequest struct {
	AppCode      string `json:"app_code" validate:"required"`
	MerchantCode string `json:"merchant_code" validate:"required"`
	TerminalCode string `json:"terminal_code" validate:"required"`
	MID          string `json:"mid"`
	SerialNumber string `json:"serial_number"`
	VaIpaymu     string `json:"va_ipaymu"`
}

type TsmGenerateApplinkRequest struct {
	AppCode      string  `json:"app_code"`
	Amount       float64 `json:"amount"`
	TrxID        string  `json:"trx_id"`
	TerminalCode string  `json:"terminal_code"`
	MerchantCode string  `json:"merchant_code"`
}

type TsmCallbackRequest struct {
	PartnerTrxID     string  `json:"partner_trx_id"`
	MerchantCode     string  `json:"merchant_code"`
	AppCode          string  `json:"app_code"`
	TerminalCode     string  `json:"terminal_code"`
	Amount           float64 `json:"amount"`
	IssuerName       string  `json:"issuer_name"`
	AcquirerHostType string  `json:"acquirer_host_type"`
	Status           string  `json:"status"`
	ResponseCode     string  `json:"response_code"`
	ResponseMessage  string  `json:"response_message"`
	DetailMessage    string  `json:"detail_message"`
}
