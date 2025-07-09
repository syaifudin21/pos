package services

import (
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type StockMovementService struct {
	DB *gorm.DB
}

func NewStockMovementService(db *gorm.DB) *StockMovementService {
	return &StockMovementService{DB: db}
}

func (s *StockMovementService) CreateStockMovement(movement *models.StockMovement) error {
	return s.DB.Create(movement).Error
}

func (s *StockMovementService) CreateStockMovementWithTx(tx *gorm.DB, movement *models.StockMovement) error {
	return tx.Create(movement).Error
}
