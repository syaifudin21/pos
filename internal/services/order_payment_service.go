package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type OrderPaymentService struct {
	DB                 *gorm.DB
	UserContextService *UserContextService
	IpaymuService      *IpaymuService
}

func NewOrderPaymentService(db *gorm.DB, userContextService *UserContextService, ipaymuService *IpaymuService) *OrderPaymentService {
	return &OrderPaymentService{DB: db, UserContextService: userContextService, IpaymuService: ipaymuService}
}

func (s *OrderPaymentService) CreateOrderPayment(req dtos.CreateOrderPaymentRequest, userID uint) (*dtos.OrderPaymentResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var order models.Order
	if err := s.DB.Where("uuid = ? AND user_id = ?", req.OrderUuid, ownerID).First(&order).Error; err != nil {
		return nil, errors.New("order not found")
	}

	if order.Status == "completed" {
		return nil, errors.New("order is already completed")
	}

	var paymentMethod models.PaymentMethod
	if err := s.DB.Where("id = ? AND is_active = ?", req.PaymentMethodID, true).First(&paymentMethod).Error; err != nil {
		return nil, errors.New("payment method not found or not active")
	}

	// Conditional validation for iPaymu
	if paymentMethod.Issuer == "iPaymu" {
		if req.CustomerName == "" {
			return nil, errors.New("customer name is required for iPaymu payments")
		}
		if req.CustomerEmail == "" {
			return nil, errors.New("customer email is required for iPaymu payments")
		}
		if req.CustomerPhone == "" {
			return nil, errors.New("customer phone is required for iPaymu payments")
		}
	}

	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	amountToPay := req.AmountPaid
	changeAmount := 0.0

	// Calculate remaining amount to be paid
	remainingAmount := order.TotalAmount - order.PaidAmount

	if amountToPay > remainingAmount {
		changeAmount = amountToPay - remainingAmount
		amountToPay = remainingAmount // Only pay the remaining amount
	}

	now := time.Now()

	orderPayment := models.OrderPayment{
		OrderID:         order.ID,
		PaymentMethodID: req.PaymentMethodID,
		AmountPaid:      amountToPay,
		PaidAt:          nil,   // Set to nil initially
		IsPaid:          false, // Initially set to false
		ChangeAmount:    changeAmount,
		CustomerName:    req.CustomerName,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
	}

	// Handle non-iPaymu payments first
	if paymentMethod.Issuer != "iPaymu" {
		orderPayment.IsPaid = true
		orderPayment.PaidAt = &now

		// Update order's paid amount
		order.PaidAmount += amountToPay
		if order.PaidAmount >= order.TotalAmount {
			order.Status = "completed"
		}

		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			log.Printf("Error updating order paid amount for non-iPaymu: %v", err)
			return nil, errors.New("failed to update order paid amount")
		}
	}

	// Ensure Extra is an empty JSON object if not set by iPaymu
	if orderPayment.Extra == "" {
		orderPayment.Extra = "{}"
	}

	// Create the order payment record for all payment types
	if err := tx.Create(&orderPayment).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating order payment: %v", err)
		return nil, errors.New("failed to create order payment")
	}

	// Handle iPaymu payment AFTER creating the base order payment record
	if paymentMethod.Issuer == "iPaymu" {
		// Preload order items to get product details
		if len(order.OrderItems) == 0 {
			if err := s.DB.Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").First(&order, order.ID).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("failed to preload order items for iPaymu payment")
			}
		}

		var products []string
		var qtys []int
		var prices []int
		for _, item := range order.OrderItems {
			products = append(products, item.ProductName)
			qtys = append(qtys, int(item.Quantity))
			prices = append(prices, int(item.Price))
		}

		log.Printf("Calling iPaymuService.CreateDirectPayment with products: %v, qtys: %v, prices: %v", products, qtys, prices)
		ipaymuRes, err := s.IpaymuService.CreateDirectPayment(
			userID,
			"Order Payment",            // ServiceName
			orderPayment.Uuid.String(), // ServiceRefID
			products,
			qtys,
			prices,
			req.CustomerName,
			req.CustomerEmail,
			req.CustomerPhone,
			paymentMethod.PaymentMethod,
			paymentMethod.PaymentChannel,
			nil, // account (optional)
		)
		if err != nil {
			tx.Rollback()
			log.Printf("iPaymu direct payment failed: %v", err)
			return nil, fmt.Errorf("iPaymu direct payment failed: %w", err)
		}

		if data, ok := ipaymuRes["Data"].(map[string]interface{}); ok {
			if trxId, ok := data["TransactionId"].(string); ok {
				orderPayment.ReferenceID = trxId
			}
			// Store the full iPaymu response in the Extra field
			rawExtra, err := json.Marshal(data)
			if err != nil {
				log.Printf("Failed to marshal iPaymu response to JSON for Extra field: %v", err)
				// Continue without extra data if marshaling fails
			} else {
				orderPayment.Extra = string(rawExtra)
			}
		}

		// Update the order payment record with iPaymu details
		if err := tx.Save(&orderPayment).Error; err != nil {
			tx.Rollback()
			log.Printf("Error updating order payment with iPaymu details: %v", err)
			return nil, errors.New("failed to update order payment with iPaymu details")
		}
	}

	// Only commit if iPaymu payment or if order status was updated for non-iPaymu
	// For iPaymu, the order status update will happen in the callback
	// Always commit the transaction here, as the orderPayment record is created.
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Error committing order payment transaction: %v", err)
		return nil, errors.New("failed to complete order payment")
	}

	// Unmarshal Extra field from string to interface{} for the response DTO
	var extraData interface{}
	if orderPayment.Extra != "" {
		if err := json.Unmarshal([]byte(orderPayment.Extra), &extraData); err != nil {
			log.Printf("Failed to unmarshal Extra field for response DTO: %v", err)
			extraData = nil // Set to nil if unmarshaling fails
		}
	}

	return &dtos.OrderPaymentResponse{
		Uuid:            orderPayment.Uuid,
		OrderUuid:       order.Uuid,
		PaymentMethodID: orderPayment.PaymentMethodID,
		PaymentName:     paymentMethod.Name,
		AmountPaid:      orderPayment.AmountPaid,
		CustomerName:    orderPayment.CustomerName,
		CustomerEmail:   orderPayment.CustomerEmail,
		CustomerPhone:   orderPayment.CustomerPhone,
		CreatedAt:       orderPayment.CreatedAt.Format("2006-01-02 15:04:05"),
		IsPaid:          orderPayment.IsPaid,
		PaidAt:          orderPayment.PaidAt,
		ChangeAmount:    orderPayment.ChangeAmount,
		Extra:           extraData,
	}, nil
}

// updateOrderAndPaymentStatus is a helper function to update order payment and order status
func (s *OrderPaymentService) updateOrderAndPaymentStatus(tx *gorm.DB, orderPayment *models.OrderPayment, amountPaid float64) error {
	now := time.Now()
	orderPayment.IsPaid = true
	orderPayment.PaidAt = &now

	if err := tx.Save(orderPayment).Error; err != nil {
		return fmt.Errorf("failed to update order payment status: %w", err)
	}

	var order models.Order
	if err := tx.Where("id = ?", orderPayment.OrderID).First(&order).Error; err != nil {
		return fmt.Errorf("order not found for order payment %s: %w", orderPayment.Uuid.String(), err)
	}

	// Update order's paid amount and status
	order.PaidAmount += amountPaid
	if order.PaidAmount >= order.TotalAmount {
		order.Status = "completed"
	}

	if err := tx.Save(&order).Error; err != nil {
		return fmt.Errorf("failed to update order paid amount and status: %w", err)
	}

	return nil
}

// UpdateOrderPaymentAndStatus updates the order payment and order status based on iPaymu notification
func (s *OrderPaymentService) UpdateOrderPaymentAndStatus(tx *gorm.DB, serviceRefID string, amountPaid float64) error {
	var orderPayment models.OrderPayment
	if err := tx.Where("uuid = ?", serviceRefID).First(&orderPayment).Error; err != nil {
		return fmt.Errorf("order payment not found for ref ID %s: %w", serviceRefID, err)
	}

	return s.updateOrderAndPaymentStatus(tx, &orderPayment, amountPaid)
}
