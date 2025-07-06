package models

import "time"

// IpaymuLog is used to log all requests and responses to/from Ipaymu
// including payment status and important timestamps.
type IpaymuLog struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ServiceName     string     `json:"service_name"`   // Nama service yang melakukan pembayaran (misal: billing, order, dsb)
	ServiceRefID    string     `json:"service_ref_id"` // ID referensi dari service terkait (misal: billing_id, order_id, dsb)
	ReferenceIpaymu string     `json:"reference_ipaymu"`
	Amount          float64    `json:"amount"`                        // Nominal pembayaran
	Status          string     `gorm:"default:pending" json:"status"` // Status pembayaran (misal: pending, paid, failed)
	PaymentMethod   string     `json:"payment_method"`                // Metode pembayaran (misal: va, qris, dsb)
	PaymentChannel  string     `json:"payment_channel"`               // Channel pembayaran (misal: bca, mandiri, dsb)
	RequestAt       time.Time  `json:"request_at"`
	SuccessAt       *time.Time `json:"success_at"`
	SettlementAt    *time.Time `json:"settlement_at"`
}
