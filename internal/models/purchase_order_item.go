package models

import (
	"github.com/google/uuid"
)

type PurchaseOrderItem struct {
	BaseModel
	PurchaseOrderID   uint            `gorm:"not null" json:"purchase_order_id"`
	PurchaseOrderUuid uuid.UUID       `gorm:"type:uuid;not null" json:"purchase_order_uuid"`
	ProductID         *uint           `gorm:"index" json:"product_id,omitempty"`
	Product           *Product        `json:"product,omitempty"`
	ProductVariantID  *uint           `gorm:"index" json:"product_variant_id,omitempty"`
	ProductVariant    *ProductVariant `json:"product_variant,omitempty"`
	Quantity          float64         `gorm:"not null" json:"quantity"`
	Price             float64         `gorm:"not null" json:"price"` // Price at the time of PO
}
