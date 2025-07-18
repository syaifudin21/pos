package services

import (
	"encoding/json"
	"errors"
	"fmt"
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

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var order models.Order
	if err := tx.Preload("OrderItems").Where("uuid = ? AND user_id = ?", req.OrderUuid, ownerID).First(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("order not found")
	}

	if order.Status == "completed" {
		tx.Rollback()
		return nil, errors.New("order is already completed")
	}

	var paymentMethod models.PaymentMethod
	if err := tx.Where("id = ? AND is_active = ?", req.PaymentMethodID, true).First(&paymentMethod).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("payment method not found or not active")
	}

	if paymentMethod.Issuer == "iPaymu" {
		if req.CustomerName == "" || req.CustomerEmail == "" || req.CustomerPhone == "" {
			tx.Rollback()
			return nil, errors.New("customer details are required for iPaymu payments")
		}
	}

	var selectedOrderItems []models.OrderItem
	if err := tx.Where("id IN ? AND order_id = ?", req.OrderItemIDs, order.ID).Find(&selectedOrderItems).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("error fetching selected order items")
	}

	if len(selectedOrderItems) != len(req.OrderItemIDs) {
		tx.Rollback()
		return nil, errors.New("one or more order items not found or do not belong to this order")
	}

	var totalAmountToPay float64 = 0
	var alreadyPaidItems []uint

	for _, item := range selectedOrderItems {
		var totalPaidQuantity float64
		// Check how much of this item has been paid for in previous transactions
		err := tx.Model(&models.OrderPaymentItem{}).
			Joins("JOIN order_payments ON order_payments.id = order_payment_items.order_payment_id").
			Where("order_payment_items.order_item_id = ? AND order_payments.is_paid = ?", item.ID, true).
			Select("COALESCE(SUM(order_payment_items.quantity_paid), 0)").
			Row().
			Scan(&totalPaidQuantity)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, fmt.Errorf("failed to check payment status for item %d: %w", item.ID, err)
		}

		if totalPaidQuantity >= item.Quantity {
			alreadyPaidItems = append(alreadyPaidItems, item.ID)
		} else {
			// For now, assume we pay the full remaining quantity. Partial payment of a single item is not supported yet.
			totalAmountToPay += item.Price * item.Quantity
		}
	}

	if len(alreadyPaidItems) > 0 {
		tx.Rollback()
		return nil, fmt.Errorf("the following items have already been fully paid: %v", alreadyPaidItems)
	}

	if totalAmountToPay <= 0 {
		tx.Rollback()
		return nil, errors.New("no amount to pay for the selected items")
	}

	now := time.Now()
	orderPayment := models.OrderPayment{
		OrderID:         order.ID,
		PaymentMethodID: req.PaymentMethodID,
		AmountPaid:      totalAmountToPay,
		IsPaid:          false,
		CustomerName:    req.CustomerName,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
		ChangeAmount:    0,
		Extra:           "{}",
	}

	var paymentItems []models.OrderPaymentItem
	for _, item := range selectedOrderItems {
		paymentItems = append(paymentItems, models.OrderPaymentItem{
			OrderItemID:  item.ID,
			QuantityPaid: item.Quantity,
		})
	}
	orderPayment.OrderPaymentItems = paymentItems

	if paymentMethod.Issuer != "iPaymu" {
		orderPayment.IsPaid = true
		orderPayment.PaidAt = &now
		order.PaidAmount += totalAmountToPay
		if order.PaidAmount >= order.TotalAmount {
			order.Status = "completed"
		}
		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update order status")
		}
	}

	if err := tx.Create(&orderPayment).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create order payment")
	}

	if paymentMethod.Issuer == "iPaymu" {
		var products []string
		var qtys []int
		var prices []int
		for _, item := range selectedOrderItems {
			products = append(products, item.ProductName)
			qtys = append(qtys, int(item.Quantity))
			prices = append(prices, int(item.Price))
		}

		ipaymuRes, err := s.IpaymuService.CreateDirectPayment(
			userID, "Order Payment", orderPayment.Uuid.String(),
			products, qtys, prices,
			req.CustomerName, req.CustomerEmail, req.CustomerPhone,
			paymentMethod.PaymentMethod, paymentMethod.PaymentChannel, nil,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("iPaymu direct payment failed: %w", err)
		}

		if data, ok := ipaymuRes["Data"].(map[string]interface{}); ok {
			if trxId, ok := data["TransactionId"].(string); ok {
				orderPayment.ReferenceID = trxId
			}
			rawExtra, err := json.Marshal(data)
			if err == nil {
				orderPayment.Extra = string(rawExtra)
			}
		}
		if err := tx.Save(&orderPayment).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update order payment with iPaymu details")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	var extraData interface{}
	if orderPayment.Extra != "" {
		json.Unmarshal([]byte(orderPayment.Extra), &extraData)
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
