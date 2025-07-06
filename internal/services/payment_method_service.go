package services

import (
	"errors"

	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type PaymentMethodService struct {
	DB *gorm.DB
}

func NewPaymentMethodService(db *gorm.DB) *PaymentMethodService {
	return &PaymentMethodService{DB: db}
}

func (s *PaymentMethodService) CreatePaymentMethod(name, payType string, creatorID uint) (*models.PaymentMethod, error) {
	paymentMethod := models.PaymentMethod{
		Name:      name,
		Type:      payType,
		IsActive:  true,
		CreatorID: &creatorID,
	}

	if err := s.DB.Create(&paymentMethod).Error; err != nil {
		return nil, err
	}

	return &paymentMethod, nil
}

func (s *PaymentMethodService) GetPaymentMethods(creatorID uint) ([]models.PaymentMethod, error) {
	var paymentMethods []models.PaymentMethod
	if err := s.DB.Where("creator_id = ?", creatorID).Find(&paymentMethods).Error; err != nil {
		return nil, err
	}
	return paymentMethods, nil
}

func (s *PaymentMethodService) UpdatePaymentMethod(id uint, name, payType string, isActive bool) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	if err := s.DB.First(&paymentMethod, id).Error; err != nil {
		return nil, errors.New("payment method not found")
	}

	paymentMethod.Name = name
	paymentMethod.Type = payType
	paymentMethod.IsActive = isActive

	if err := s.DB.Save(&paymentMethod).Error; err != nil {
		return nil, err
	}

	return &paymentMethod, nil
}

func (s *PaymentMethodService) DeletePaymentMethod(id uint) error {
	var paymentMethod models.PaymentMethod
	if err := s.DB.First(&paymentMethod, id).Error; err != nil {
		return errors.New("payment method not found")
	}

	return s.DB.Delete(&paymentMethod).Error
}
