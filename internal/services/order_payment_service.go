package services

import (
	"errors"
	"log"
	"time"

	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type OrderPaymentService struct {
	DB                 *gorm.DB
	UserContextService *UserContextService
}

func NewOrderPaymentService(db *gorm.DB, userContextService *UserContextService) *OrderPaymentService {
	return &OrderPaymentService{DB: db, UserContextService: userContextService}
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
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		PaidAt:          &now,
		IsPaid:          true,
		ChangeAmount:    changeAmount,
	}

	if err := tx.Create(&orderPayment).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating order payment: %v", err)
		return nil, errors.New("failed to create order payment")
	}

	// Update order's paid amount
	order.PaidAmount += amountToPay
	if order.PaidAmount >= order.TotalAmount {
		order.Status = "completed"
	}

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating order paid amount: %v", err)
		return nil, errors.New("failed to update order paid amount")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Error committing order payment transaction: %v", err)
		return nil, errors.New("failed to complete order payment")
	}

	return &dtos.OrderPaymentResponse{
		Uuid:            orderPayment.Uuid,
		OrderUuid:       order.Uuid,
		PaymentMethodID: orderPayment.PaymentMethodID,
		PaymentName:     paymentMethod.Name,
		AmountPaid:      orderPayment.AmountPaid,
		CustomerName:    orderPayment.CustomerName,
		CustomerPhone:   orderPayment.CustomerPhone,
		CreatedAt:       orderPayment.CreatedAt.Format("2006-01-02 15:04:05"),
		IsPaid:          orderPayment.IsPaid,
		PaidAt:          orderPayment.PaidAt,
		ChangeAmount:    orderPayment.ChangeAmount,
	}, nil
}
