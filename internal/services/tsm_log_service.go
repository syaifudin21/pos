package services

import (
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type TsmLogService struct {
	DB *gorm.DB
}

func NewTsmLogService(db *gorm.DB) *TsmLogService {
	return &TsmLogService{
		DB: db,
	}
}

func (s *TsmLogService) CreateTsmLog(logEntry *models.TsmLog) error {
	return s.DB.Create(logEntry).Error
}
