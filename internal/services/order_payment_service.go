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
		PaidAt:          &now,
		IsPaid:          false, // Initially set to false
		ChangeAmount:    changeAmount,
		CustomerName:    req.CustomerName,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
	}

	// Handle iPaymu payment
	if paymentMethod.Issuer == "iPaymu" {
		// For iPaymu, IsPaid remains false until callback
		// Order status and paid amount will be updated in iPaymu callback
		// Preload order items to get product details
		// Note: This preloading is done outside the main transaction to avoid issues if the transaction rolls back.
		// However, for iPaymu, we need the product details before creating the payment record.
		// If order items are not already loaded, load them here.
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
			"Order Payment",     // ServiceName
			order.Uuid.String(), // ServiceRefID
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
		}
		// Store the full iPaymu response in the Extra field
		rawExtra, err := json.Marshal(ipaymuRes["Data"].(map[string]interface{}))
		if err != nil {
			log.Printf("Failed to marshal iPaymu response to JSON for Extra field: %v", err)
			// Continue without extra data if marshaling fails
		} else {
			orderPayment.Extra = string(rawExtra)
		}
		// You might want to update PaidAt based on iPaymu response if it provides a specific timestamp
	} else {
		// For non-iPaymu payments, mark as paid and update order status immediately
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

	if err := tx.Create(&orderPayment).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating order payment: %v", err)
		return nil, errors.New("failed to create order payment")
	}

	// Only commit if iPaymu payment or if order status was updated for non-iPaymu
	// For iPaymu, the order status update will happen in the callback
	if paymentMethod.Issuer == "iPaymu" || orderPayment.IsPaid {
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Printf("Error committing order payment transaction: %v", err)
			return nil, errors.New("failed to complete order payment")
		}
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
