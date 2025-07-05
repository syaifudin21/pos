package models

import (
	"github.com/google/uuid"
)

type PurchaseOrder struct {
	BaseModel
	SupplierID         uint                `gorm:"not null" json:"supplier_id"`
	Supplier           Supplier            `json:"supplier"`
	OutletID           uint                `gorm:"not null" json:"outlet_id"`
	Outlet             Outlet              `json:"outlet"`
	Status             string              `gorm:"not null" json:"status"` // e.g., "pending", "completed", "cancelled"
	TotalAmount        float64             `gorm:"not null" json:"total_amount"`
	PurchaseOrderItems []PurchaseOrderItem `json:"purchase_order_items"`
	UserID             uuid.UUID           `gorm:"type:uuid;not null" json:"user_id"`
	User               User                `json:"user"`
}
