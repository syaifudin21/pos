package dtos

type CreateDirectPaymentRequest struct {
	ServiceName  string   `json:"service_name" validate:"required"`
	ServiceRefID string   `json:"service_ref_id" validate:"required,uuid"` // Order UUID
	Product      []string `json:"product" validate:"required,min=1,dive,required"`
	Qty          []int    `json:"qty" validate:"required,min=1,dive,gt=0"`
	Price        []int    `json:"price" validate:"required,min=1,dive,gt=0"`
	Name         string   `json:"name" validate:"required"`
	Email        string   `json:"email" validate:"required,email"`
	Phone        string   `json:"phone" validate:"required"`
	Method       string   `json:"method" validate:"required"`
	Channel      string   `json:"channel" validate:"required"`
	Account      *string  `json:"account,omitempty"`
}

type IpaymuResponse struct {
	Status  int                  `json:"Status"`
	Message string               `json:"Message"`
	Success bool                 `json:"Success"`
	Data    IpaymuResponseDetail `json:"Data"`
}

type IpaymuResponseDetail struct {
	Channel       string  `json:"Channel"`
	Escrow        bool    `json:"Escrow"`
	Expired       string  `json:"Expired"`
	Fee           int     `json:"Fee"`
	FeeDirection  string  `json:"FeeDirection"`
	Note          *string `json:"Note"`
	PaymentName   string  `json:"PaymentName"`
	PaymentNo     string  `json:"PaymentNo"`
	ReferenceId   int     `json:"ReferenceId"`
	SessionId     int     `json:"SessionId"`
	SubTotal      int     `json:"SubTotal"`
	Total         int     `json:"Total"`
	TransactionId int     `json:"TransactionId"`
	Via           string  `json:"Via"`
}

type IpaymuNotifyRequest struct {
	TrxID                 int           `json:"trx_id"`                  // 170981
	SID                   string        `json:"sid"`                     // "1"
	ReferenceID           string        `json:"reference_id"`            // "1"
	Status                string        `json:"status"`                  // "berhasil"
	StatusCode            int           `json:"status_code"`             // 1
	SubTotal              string        `json:"sub_total"`               // "130000"
	Total                 string        `json:"total"`                   // "134000"
	Amount                string        `json:"amount"`                  // "134000"
	Fee                   string        `json:"fee"`                     // "4000"
	PaidOff               int           `json:"paid_off"`                // 130000
	CreatedAt             string        `json:"created_at"`              // "2025-07-06 19:01:29"
	ExpiredAt             string        `json:"expired_at"`              // "2025-07-07 19:01:29"
	PaidAt                string        `json:"paid_at"`                 // "2025-07-06 19:01:45"
	SettlementStatus      string        `json:"settlement_status"`       // "settled"
	TransactionStatusCode int           `json:"transaction_status_code"` // 1
	IsEscrow              bool          `json:"is_escrow"`               // false
	SystemNotes           string        `json:"system_notes"`            // "Sandbox notify"
	Via                   string        `json:"via"`                     // "va"
	Channel               string        `json:"channel"`                 // "mandiri"
	PaymentNo             string        `json:"payment_no"`              // "000012316415"
	BuyerName             string        `json:"buyer_name"`              // "Budi Santoso"
	BuyerEmail            string        `json:"buyer_email"`             // "budi@example.com"
	BuyerPhone            string        `json:"buyer_phone"`             // "08123456789"
	AdditionalInfo        []interface{} `json:"additional_info"`         // []
	URL                   string        `json:"url"`                     // "http://localhost:8080/vendor/ipaymu/notify"
	VA                    string        `json:"va"`                      // "000012316415"
}

// RegisterIpaymuRequest untuk pendaftaran user ke Ipaymu
// Untuk identityPhoto, gunakan multipart/form-data jika ingin upload file
// Field validate opsional, bisa disesuaikan dengan kebutuhan
// Email opsional, jika tidak ada, withoutEmail harus '1'
type RegisterIpaymuRequest struct {
	Name         string  `json:"name" validate:"required"`
	Phone        string  `json:"phone" validate:"required"`
	Password     string  `json:"password" validate:"required"`
	Email        *string `json:"email,omitempty"`
	WithoutEmail string  `json:"withoutEmail"`
	IdentityNo   *string `json:"identityNo,omitempty"`
	BusinessName *string `json:"businessName,omitempty"`
	Birthday     *string `json:"birthday,omitempty"`
	Birthplace   *string `json:"birthplace,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	Address      *string `json:"address,omitempty"`
	// Untuk multipart/form-data, gunakan field ini untuk file
	IdentityPhoto interface{} `json:"identityPhoto,omitempty"`
}
