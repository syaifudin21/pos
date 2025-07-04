package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type OrderService struct {
	DB           *gorm.DB
	StockService *StockService // Dependency on StockService
}

func NewOrderService(db *gorm.DB, stockService *StockService) *OrderService {
	return &OrderService{DB: db, StockService: stockService}
}

// CreateOrder creates a new order and deducts stock.
func (s *OrderService) CreateOrder(outletUuid, userUuid uuid.UUID, items []models.OrderItemRequest) (*models.Order, error) {
	// Find Outlet
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", outletUuid).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	// Find User
	var user models.User
	if err := s.DB.Where("uuid = ?", userUuid).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Start a database transaction
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	order := models.Order{
		OutletID:    outlet.ID,
		UserID:      user.ID,
		Status:      "completed", // Assuming immediate completion for simplicity
		TotalAmount: 0,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating order: %v", err)
		return nil, errors.New("failed to create order")
	}

	totalAmount := 0.0
	for _, item := range items {
		var product models.Product
		if err := tx.Where("uuid = ?", item.ProductUuid).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}

		// Deduct stock using StockService
		if err := s.StockService.DeductStockForSale(outletUuid, item.ProductUuid, item.Quantity); err != nil {
			tx.Rollback()
			return nil, err // Return specific stock deduction error
		}

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price, // Price at the time of order
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating order item: %v", err)
			return nil, errors.New("failed to create order item")
		}
		totalAmount += product.Price * item.Quantity
	}

	order.TotalAmount = totalAmount
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating order total amount: %v", err)
		return nil, errors.New("failed to update order total amount")
	}

	tx.Commit()

	// Reload order with associations for response
	s.DB.Preload("Outlet").Preload("User").Preload("OrderItems.Product").First(&order, order.ID)
	return &order, nil
}

// GetOrder retrieves an order by its external ID.
func (s *OrderService) GetOrderByUuid(externalID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := s.DB.Preload("Outlet").Preload("User").Preload("OrderItems.Product").Where("uuid = ?", externalID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		log.Printf("Error getting order by ExternalID: %v", err)
		return nil, errors.New("failed to retrieve order")
	}
	return &order, nil
}

// GetOrdersByOutlet retrieves all orders for a specific outlet.
func (s *OrderService) GetOrdersByOutlet(outletUuid uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", outletUuid).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	if err := s.DB.Preload("User").Preload("OrderItems.Product").Where("outlet_id = ?", outlet.ID).Find(&orders).Error; err != nil {
		log.Printf("Error getting orders by outlet: %v", err)
		return nil, errors.New("failed to retrieve orders")
	}
	return orders, nil
}
