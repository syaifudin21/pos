package dtos

type IpaymuRequest struct {
	Product    []string `json:"product"`
	Qty        []int    `json:"qty"`
	Price      []int    `json:"price"`
	ReturnUrl  string   `json:"returnUrl"`
	CancelUrl  string   `json:"cancelUrl"`
	NotifyUrl  string   `json:"notifyUrl"`
	ReferenceId string   `json:"referenceId"`
	BuyerName  string   `json:"buyerName"`
	BuyerEmail string   `json:"buyerEmail"`
	BuyerPhone string   `json:"buyerPhone"`
	Udf1       string   `json:"udf1"` // Custom field for order UUID
}

type IpaymuResponse struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Data     IpaymuResponseData `json:"data"`
	Comments string `json:"comments"`
}

type IpaymuResponseData struct {
	SessionID    string `json:"sessionId"`
	URL          string `json:"url"`
	TransactionID string `json:"transactionId"`
}
