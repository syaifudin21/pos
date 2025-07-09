package models

import (
	"time"

	"gorm.io/gorm"
)

// StockMovement represents a record of stock changes
type StockMovement struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ProductID     uint           `gorm:"not null"`
	Product       Product        `gorm:"foreignKey:ProductID"`
	OutletID      uint           `gorm:"not null"`
	Outlet        Outlet         `gorm:"foreignKey:OutletID"`
	QuantityChange int            `gorm:"not null"` // Positive for increase, negative for decrease
	MovementType  string         `gorm:"type:varchar(50);not null"` // e.g., "Order", "PurchaseOrder", "Adjustment"
	ReferenceID   *uint          // Optional: ID of the Order, PurchaseOrder, or other reference
	Description   *string        // Optional: A brief description for manual adjustments
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
