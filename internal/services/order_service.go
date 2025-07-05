package services

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type OrderService struct {
	DB           *gorm.DB
	StockService *StockService // Dependency on StockService
	IpaymuService *IpaymuService // Dependency on IpaymuService
}

func NewOrderService(db *gorm.DB, stockService *StockService, ipaymuService *IpaymuService) *OrderService {
	return &OrderService{DB: db, StockService: stockService, IpaymuService: ipaymuService}
}

// GetOwnerID retrieves the owner's ID for a given user.
// If the user is a manager or cashier, it returns their creator's ID.
// Otherwise, it returns the user's own ID.
func (s *OrderService) GetOwnerID(userID uint) (uint, error) {
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

// CreateOrder creates a new order and deducts stock.
func (s *OrderService) CreateOrder(outletUuid uuid.UUID, userID uint, items []dtos.OrderItemRequest, paymentMethod string) (*dtos.OrderResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	// Find Outlet
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	// Find User
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Start a database transaction
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	order := models.Order{
		OutletID:      outlet.ID,
		UserID:        ownerID,
		Status:        "completed", // Assuming immediate completion for simplicity
		TotalAmount:   0,
		PaymentMethod: paymentMethod,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating order: %v", err)
		return nil, errors.New("failed to create order")
	}

	totalAmount := 0.0
	for _, item := range items {
		var product models.Product
		if err := tx.Where("uuid = ? AND user_id = ?", item.ProductUuid, ownerID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}

		// Deduct stock using StockService
		if err := s.StockService.DeductStockForSale(outletUuid, item.ProductUuid, float64(item.Quantity), ownerID); err != nil {
			tx.Rollback()
			return nil, err // Return specific stock deduction error
		}

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  float64(item.Quantity),
			Price:     product.Price, // Price at the time of order
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating order item: %v", err)
			return nil, errors.New("failed to create order item")
		}
		totalAmount += product.Price * float64(item.Quantity)
	}

	order.TotalAmount = totalAmount
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating order total amount: %v", err)
		return nil, errors.New("failed to update order total amount")
	}

	// Process payment if method is iPaymu
	if paymentMethod == "ipaymu" {
		// Fetch order items with product details for iPaymu
		var fullOrderItems []models.OrderItem
		if err := tx.Preload("Product").Where("order_id = ?", order.ID).Find(&fullOrderItems).Error; err != nil {
			tx.Rollback()
			log.Printf("Error preloading order items for iPaymu: %v", err)
			return nil, errors.New("failed to prepare order for iPaymu payment")
		}

		ipaymuResp, err := s.IpaymuService.ProcessIpaymuPayment(&order, &user, fullOrderItems)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// You might want to store ipaymuResp.Data.URL or transaction ID in the order for later reference
		log.Printf("iPaymu Payment URL: %s", ipaymuResp.Data.URL)
	}

	tx.Commit()

	// Construct DTO response
	return &dtos.OrderResponse{
		ID:            order.ID,
		Uuid:          order.Uuid,
		OutletID:      order.OutletID,
		OutletUuid:    outlet.Uuid,
		UserID:        order.UserID,
		UserUuid:      user.Uuid,
		OrderDate:     order.CreatedAt.Format(time.RFC3339),
		TotalAmount:   order.TotalAmount,
		PaymentMethod: order.PaymentMethod,
		Status:        order.Status,
	}, nil
}

// GetOrder retrieves an order by its Uuid.
func (s *OrderService) GetOrderByUuid(uuid uuid.UUID, userID uint) (*dtos.OrderResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var order models.Order
	if err := s.DB.Preload("Outlet").Preload("User").Where("uuid = ? AND user_id = ?", uuid, ownerID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		log.Printf("Error getting order by uuid: %v", err)
		return nil, errors.New("failed to retrieve order")
	}

	return &dtos.OrderResponse{
		ID:            order.ID,
		Uuid:          order.Uuid,
		OutletID:      order.OutletID,
		OutletUuid:    order.Outlet.Uuid,
		UserID:        order.UserID,
		UserUuid:      order.User.Uuid,
		OrderDate:     order.CreatedAt.Format(time.RFC3339),
		TotalAmount:   order.TotalAmount,
		PaymentMethod: order.PaymentMethod,
		Status:        order.Status,
	}, nil
}

// GetOrdersByOutlet retrieves all orders for a specific outlet.
func (s *OrderService) GetOrdersByOutlet(outletUuid uuid.UUID, userID uint) ([]dtos.OrderResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	if err := s.DB.Preload("User").Where("outlet_id = ? AND user_id = ?", outlet.ID, ownerID).Find(&orders).Error; err != nil {
		log.Printf("Error getting orders by outlet: %v", err)
		return nil, errors.New("failed to retrieve orders")
	}

	var orderResponses []dtos.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, dtos.OrderResponse{
			ID:            order.ID,
			Uuid:          order.Uuid,
			OutletID:      order.OutletID,
			OutletUuid:    outlet.Uuid, // Use the fetched outlet's UUID
			UserID:        order.UserID,
			UserUuid:      order.User.Uuid,
			OrderDate:     order.CreatedAt.Format(time.RFC3339),
			TotalAmount:   order.TotalAmount,
			PaymentMethod: order.PaymentMethod,
			Status:        order.Status,
		})
	}
	return orderResponses, nil
}
