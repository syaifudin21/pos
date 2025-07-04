package models

import (
	"github.com/google/uuid"
)

type PurchaseOrderItem struct {
	BaseModel
	PurchaseOrderID   uint      `gorm:"not null" json:"purchase_order_id"`
	PurchaseOrderUuid uuid.UUID `gorm:"type:uuid;not null" json:"purchase_order_uuid"`
	ProductID         uint      `gorm:"not null" json:"product_id"`
	Product           Product   `json:"product"`
	Quantity          float64   `gorm:"not null" json:"quantity"`
	Price             float64   `gorm:"not null" json:"price"` // Price at the time of PO
}
