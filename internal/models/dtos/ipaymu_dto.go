package dtos

type CreateDirectPaymentRequest struct {
	Product  []string `json:"product" validate:"required,min=1,dive,required"`
	Qty      []int    `json:"qty" validate:"required,min=1,dive,gt=0"`
	Price    []int    `json:"price" validate:"required,min=1,dive,gt=0"`
	Name     string   `json:"name" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Phone    string   `json:"phone" validate:"required"`
	Callback string   `json:"callback" validate:"required,url"`
	Method   string   `json:"method" validate:"required"`
	Channel  string   `json:"channel" validate:"required"`
	Account  *string  `json:"account,omitempty"`
}

type IpaymuRequest struct {
	Product     []string `json:"product" validate:"required"`
	Qty         []int    `json:"qty" validate:"required"`
	Price       []int    `json:"price" validate:"required"`
	ReturnUrl   string   `json:"returnUrl"`
	CancelUrl   string   `json:"cancelUrl"`
	NotifyUrl   string   `json:"notifyUrl"`
	ReferenceId string   `json:"referenceId" validate:"required,uuid"` // Order UUID
	BuyerName   string   `json:"buyerName" validate:"required"`
	BuyerEmail  string   `json:"buyerEmail" validate:"required,email"`
	BuyerPhone  string   `json:"buyerPhone" validate:"required"`
	Udf1        string   `json:"udf1"` // Custom field for order UUID
}

type IpaymuResponse struct {
	Status   int                `json:"status"`
	Message  string             `json:"message"`
	Data     IpaymuResponseData `json:"data"`
	Comments string             `json:"comments"`
}

type IpaymuResponseData struct {
	SessionID     string `json:"sessionId"`
	URL           string `json:"url"`
	TransactionID string `json:"transactionId"`
}
