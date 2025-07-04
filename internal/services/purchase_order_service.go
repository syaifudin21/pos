package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type PurchaseOrderService struct {
	DB           *gorm.DB
	StockService *StockService // Dependency on StockService
}

func NewPurchaseOrderService(db *gorm.DB, stockService *StockService) *PurchaseOrderService {
	return &PurchaseOrderService{DB: db, StockService: stockService}
}

// CreatePurchaseOrder creates a new purchase order.
func (s *PurchaseOrderService) CreatePurchaseOrder(supplierUuid, outletUuid uuid.UUID, items []models.PurchaseOrderItemRequest) (*models.PurchaseOrder, error) {
	var supplier models.Supplier
	if err := s.DB.Where("uuid = ?", supplierUuid).First(&supplier).Error; err != nil {
		return nil, errors.New("supplier not found")
	}

	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", outletUuid).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	po := models.PurchaseOrder{
		SupplierID:  supplier.ID,
		OutletID:    outlet.ID,
		Status:      "pending",
		TotalAmount: 0,
	}

	if err := tx.Create(&po).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating purchase order: %v", err)
		return nil, errors.New("failed to create purchase order")
	}

	totalAmount := 0.0
	for _, item := range items {
		var product models.Product
		if err := tx.Where("uuid = ?", item.Productuuid).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}

		poItem := models.PurchaseOrderItem{
			PurchaseOrderID:   po.ID,
			PurchaseOrderUuid: po.Uuid, // Set the ExternalID here
			ProductID:         product.ID,
			Quantity:          item.Quantity,
			Price:             item.Price,
		}

		if err := tx.Create(&poItem).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating purchase order item: %v", err)
			return nil, errors.New("failed to create purchase order item")
		}
		totalAmount += item.Price * item.Quantity
	}

	po.TotalAmount = totalAmount
	if err := tx.Save(&po).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating purchase order total amount: %v", err)
		return nil, errors.New("failed to update purchase order total amount")
	}

	tx.Commit()

	s.DB.Preload("Supplier").Preload("Outlet").Preload("PurchaseOrderItems.Product").First(&po, po.ID)
	return &po, nil
}

// GetPurchaseOrderByExternalID retrieves a purchase order by its external ID.
func (s *PurchaseOrderService) GetPurchaseOrderByUuid(uuid uuid.UUID) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	if err := s.DB.Preload("Supplier").Preload("Outlet").Preload("PurchaseOrderItems.Product").Where("uuid = ?", uuid).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("purchase order not found")
		}
		log.Printf("Error getting purchase order by ExternalID: %v", err)
		return nil, errors.New("failed to retrieve purchase order")
	}
	return &po, nil
}

// GetPurchaseOrdersByOutlet retrieves all purchase orders for a specific outlet.
func (s *PurchaseOrderService) GetPurchaseOrdersByOutlet(outletUuid uuid.UUID) ([]models.PurchaseOrder, error) {
	var pos []models.PurchaseOrder
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", outletUuid).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	if err := s.DB.Preload("Supplier").Preload("PurchaseOrderItems.Product").Where("outlet_id = ?", outlet.ID).Find(&pos).Error; err != nil {
		log.Printf("Error getting purchase orders by outlet: %v", err)
		return nil, errors.New("failed to retrieve purchase orders")
	}
	return pos, nil
}

// ReceivePurchaseOrder updates stock based on a completed purchase order.
func (s *PurchaseOrderService) ReceivePurchaseOrder(poUuid uuid.UUID) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder

	log.Printf(po.Outlet.Name)
	if err := s.DB.Preload("PurchaseOrderItems.Product").Where("uuid = ?", poUuid).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("purchase order not found")
		}
		log.Printf("Error finding purchase order to receive: %v", err)
		return nil, errors.New("failed to retrieve purchase order")
	}

	if po.Status == "completed" {
		return nil, errors.New("purchase order already received")
	}

	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	for _, item := range po.PurchaseOrderItems {
		// Add stock using StockService
		// log.Printf("Menerima PO: OutletUuid=%v, ProductUuid=%v, Quantity=%v", po.Outlet, item.Product.Uuid, item.Quantity)

		_, err := s.StockService.UpdateStock(po.Outlet.Uuid, item.Product.Uuid, item.Quantity) // Assuming UpdateStock adds to existing quantity
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating stock for received PO: %v", err)
			return nil, errors.New("failed to update stock for received purchase order")
		}
	}

	po.Status = "completed"
	if err := tx.Save(&po).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating purchase order status: %v", err)
		return nil, errors.New("failed to update purchase order status")
	}

	tx.Commit()

	s.DB.Preload("Supplier").Preload("Outlet").Preload("PurchaseOrderItems.Product").First(&po, po.ID)
	return &po, nil
}
