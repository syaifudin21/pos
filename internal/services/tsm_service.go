package services

import (
	"errors"

	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type TsmService struct {
	DB                 *gorm.DB
	UserContextService *UserContextService
}

func NewTsmService(db *gorm.DB, userContextService *UserContextService) *TsmService {
	return &TsmService{
		DB:                 db,
		UserContextService: userContextService,
	}
}

func (s *TsmService) RegisterTsm(userID uint, req dtos.TsmRegisterRequest) error {
	var existingTsm models.UserTsm
	result := s.DB.Where("user_id = ?", userID).First(&existingTsm)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new entry
			userTsm := models.UserTsm{
				UserID:       userID,
				AppCode:      req.AppCode,
				MerchantCode: req.MerchantCode,
				TerminalCode: req.TerminalCode,
				SerialNumber: req.SerialNumber,
				MID:          req.MID,
				VaIpaymu:     req.VaIpaymu,
			}
			if err := s.DB.Create(&userTsm).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	} else {
		// Update existing entry
		existingTsm.AppCode = req.AppCode
		existingTsm.MerchantCode = req.MerchantCode
		existingTsm.TerminalCode = req.TerminalCode
		existingTsm.SerialNumber = req.SerialNumber
		existingTsm.MID = req.MID
		existingTsm.VaIpaymu = req.VaIpaymu
		if err := s.DB.Save(&existingTsm).Error; err != nil {
			return err
		}
	}

	return nil
}
