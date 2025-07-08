package services

import (
	"errors"

	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

var (
	ErrIpaymuRegistrationRequired = errors.New("iPaymu registration required")
	ErrTsmRegistrationRequired    = errors.New("TSM registration required")
)

type UserPaymentService struct {
	DB                 *gorm.DB
	UserContextService *UserContextService
}

func NewUserPaymentService(db *gorm.DB, userContextService *UserContextService) *UserPaymentService {
	return &UserPaymentService{
		DB:                 db,
		UserContextService: userContextService,
	}
}

func (s *UserPaymentService) ActivateUserPayment(userID uint, paymentMethodID uint) error {
	// Check if the payment method exists and is active
	var paymentMethod models.PaymentMethod
	if err := s.DB.Where("id = ? AND is_active = ?", paymentMethodID, true).First(&paymentMethod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment method not found or not active")
		}
		return err
	}

	// Check payment channel for special handling
	switch paymentMethod.Issuer {
	case "iPaymu":
		// Check if UserIpaymu entry exists for this user
		var userIpaymu models.UserIpaymu
		if err := s.DB.Where("user_id = ?", userID).First(&userIpaymu).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrIpaymuRegistrationRequired
			}
			return err
		}
	case "TSM":
		// Check if UserTsm entry exists for this user
		var userTsm models.UserTsm
		if err := s.DB.Where("user_id = ?", userID).First(&userTsm).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTsmRegistrationRequired
			}
			return err
		}
	}

	var userPayment models.UserPayment
	result := s.DB.Where("user_id = ? AND payment_method_id = ?", userID, paymentMethodID).First(&userPayment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new entry if not found
			userPayment = models.UserPayment{
				UserID:          userID,
				PaymentMethodID: paymentMethodID,
				IsActive:        true,
			}
			if err := s.DB.Create(&userPayment).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	} else {
		// Update existing entry
		if !userPayment.IsActive {
			userPayment.IsActive = true
			if err := s.DB.Save(&userPayment).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *UserPaymentService) DeactivateUserPayment(userID uint, paymentMethodID uint) error {
	var userPayment models.UserPayment
	result := s.DB.Where("user_id = ? AND payment_method_id = ?", userID, paymentMethodID).First(&userPayment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user payment method not found")
		}
		return result.Error
	}

	if userPayment.IsActive {
		userPayment.IsActive = false
		if err := s.DB.Save(&userPayment).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *UserPaymentService) ListUserPaymentsByOwner(ownerID uint) ([]models.UserPayment, error) {
	var userPayments []models.UserPayment

	// Get all user IDs associated with this owner (including the owner itself)
	var userIDs []uint
	userIDs = append(userIDs, ownerID)

	var subUsers []models.User
	if err := s.DB.Where("creator_id = ?", ownerID).Find(&subUsers).Error; err != nil {
		return nil, err
	}

	for _, user := range subUsers {
		userIDs = append(userIDs, user.ID)
	}

	// Find user payments for these user IDs, preloading PaymentMethod
	if err := s.DB.Preload("PaymentMethod").Where("user_id IN ?", userIDs).Find(&userPayments).Error; err != nil {
		return nil, err
	}

	return userPayments, nil
}
