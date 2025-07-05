package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type PurchaseOrderService struct {
	DB           *gorm.DB
	StockService *StockService // Dependency on StockService
}

func NewPurchaseOrderService(db *gorm.DB, stockService *StockService) *PurchaseOrderService {
	return &PurchaseOrderService{DB: db, StockService: stockService}
}

// GetOwnerID retrieves the owner's ID for a given user.
// If the user is a manager or cashier, it returns their creator's ID.
// Otherwise, it returns the user's own ID.
func (s *PurchaseOrderService) GetOwnerID(userID uint) (uint, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		log.Printf("Error finding user: %v", err)
		return 0, errors.New("user not found")
	}

	if (user.Role == "manager" || user.Role == "cashier") && user.CreatorID != nil {
		return *user.CreatorID, nil
	}

	return userID, nil
}

// CreatePurchaseOrder creates a new purchase order.
func (s *PurchaseOrderService) CreatePurchaseOrder(supplierUuid, outletUuid uuid.UUID, items []dtos.PurchaseItemRequest, userID uint) (*dtos.PurchaseOrderResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var supplier models.Supplier
	if err := s.DB.Where("uuid = ? AND user_id = ?", supplierUuid, ownerID).First(&supplier).Error; err != nil {
		return nil, errors.New("supplier not found")
	}

	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
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
		UserID:      ownerID,
	}

	if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&po).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating purchase order: %v", err)
		return nil, errors.New("failed to create purchase order")
	}

	totalAmount := 0.0
	for _, item := range items {
		var product models.Product
		if err := tx.Where("uuid = ? AND user_id = ?", item.ProductUuid, ownerID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}

		poItem := models.PurchaseOrderItem{
			PurchaseOrderID:   po.ID,
			PurchaseOrderUuid: po.Uuid, // Set the UUID here
			ProductID:         product.ID,
			Quantity:          float64(item.Quantity),
			Price:             item.Price,
		}

		if err := tx.Create(&poItem).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating purchase order item: %v", err)
			return nil, errors.New("failed to create purchase order item")
		}
		totalAmount += item.Price * float64(item.Quantity)
	}

	po.TotalAmount = totalAmount
	if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&po).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating purchase order total amount: %v", err)
		return nil, errors.New("failed to update purchase order total amount")
	}

	tx.Commit()

	return &dtos.PurchaseOrderResponse{
		ID:           po.ID,
		Uuid:         po.Uuid,
		SupplierID:   supplier.ID,
		SupplierUuid: supplier.Uuid,
		OutletID:     po.OutletID,
		OutletUuid:   outlet.Uuid,
		OrderDate:    po.CreatedAt.Format(time.RFC3339),
		TotalAmount:  po.TotalAmount,
		Status:       po.Status,
	}, nil
}

// GetPurchaseOrderByUuid retrieves a purchase order by its UUID.
func (s *PurchaseOrderService) GetPurchaseOrderByUuid(uuid uuid.UUID, userID uint) (*dtos.PurchaseOrderResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var po models.PurchaseOrder
	if err := s.DB.Preload("Supplier").Preload("Outlet").Where("uuid = ? AND user_id = ?", uuid, ownerID).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("purchase order not found")
		}
		log.Printf("Error getting purchase order by UUID: %v", err)
		return nil, errors.New("failed to retrieve purchase order")
	}
	return &dtos.PurchaseOrderResponse{
		ID:           po.ID,
		Uuid:         po.Uuid,
		SupplierID:   po.SupplierID,
		SupplierUuid: po.Supplier.Uuid,
		OutletID:     po.OutletID,
		OutletUuid:   po.Outlet.Uuid,
		OrderDate:    po.CreatedAt.Format(time.RFC3339),
		TotalAmount:  po.TotalAmount,
		Status:       po.Status,
	}, nil
}

// GetPurchaseOrdersByOutlet retrieves all purchase orders for a specific outlet.
func (s *PurchaseOrderService) GetPurchaseOrdersByOutlet(outletUuid uuid.UUID, userID uint) ([]dtos.PurchaseOrderResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var pos []models.PurchaseOrder
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	if err := s.DB.Preload("Supplier").Where("outlet_id = ? AND user_id = ?", outlet.ID, ownerID).Find(&pos).Error; err != nil {
		log.Printf("Error getting purchase orders by outlet: %v", err)
		return nil, errors.New("failed to retrieve purchase orders")
	}

	var poResponses []dtos.PurchaseOrderResponse
	for _, po := range pos {
		poResponses = append(poResponses, dtos.PurchaseOrderResponse{
			ID:           po.ID,
			Uuid:         po.Uuid,
			SupplierID:   po.SupplierID,
			SupplierUuid: po.Supplier.Uuid,
			OutletID:     po.OutletID,
			OutletUuid:   outlet.Uuid,
			OrderDate:    po.CreatedAt.Format(time.RFC3339),
			TotalAmount:  po.TotalAmount,
			Status:       po.Status,
		})
	}
	return poResponses, nil
}

// ReceivePurchaseOrder updates stock based on a completed purchase order.
func (s *PurchaseOrderService) ReceivePurchaseOrder(poUuid uuid.UUID, userID uint) (*dtos.PurchaseOrderResponse, error) {
	var po models.PurchaseOrder
	if err := s.DB.Preload("Outlet").Preload("PurchaseOrderItems.Product").Where("uuid = ? AND user_id = ?", poUuid, userID).First(&po).Error; err != nil {
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
		_, err := s.StockService.AdjustStock(po.Outlet.Uuid, item.Product.Uuid, item.Quantity, userID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating stock for received PO: %v", err)
			return nil, errors.New("failed to update stock for received purchase order")
		}
	}

	po.Status = "completed"
	if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&po).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating purchase order status: %v", err)
		return nil, errors.New("failed to update purchase order status")
	}

	tx.Commit()

	return &dtos.PurchaseOrderResponse{
		ID:           po.ID,
		Uuid:         po.Uuid,
		SupplierID:   po.SupplierID,
		SupplierUuid: po.Supplier.Uuid,
		OutletID:     po.OutletID,
		OutletUuid:   po.Outlet.Uuid,
		OrderDate:    po.CreatedAt.Format(time.RFC3339),
		TotalAmount:  po.TotalAmount,
		Status:       po.Status,
	}, nil
}
