package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type OutletService struct {
	DB *gorm.DB
}

func NewOutletService(db *gorm.DB) *OutletService {
	return &OutletService{DB: db}
}

func (s *OutletService) GetAllOutlets() ([]models.Outlet, error) {
	var outlets []models.Outlet
	if err := s.DB.Find(&outlets).Error; err != nil {
		log.Printf("Error getting all outlets: %v", err)
		return nil, errors.New("failed to retrieve outlets")
	}
	return outlets, nil
}

func (s *OutletService) GetOutletByUuid(Uuid uuid.UUID) (*models.Outlet, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", Uuid).First(&outlet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet not found")
		}
		log.Printf("Error getting outlet by Uuid: %v", err)
		return nil, errors.New("failed to retrieve outlet")
	}
	return &outlet, nil
}

func (s *OutletService) CreateOutlet(outlet *models.Outlet) (*models.Outlet, error) {
	if err := s.DB.Create(outlet).Error; err != nil {
		log.Printf("Error creating outlet: %v", err)
		return nil, errors.New("failed to create outlet")
	}
	return outlet, nil
}

func (s *OutletService) UpdateOutlet(Uuid uuid.UUID, updatedOutlet *models.Outlet) (*models.Outlet, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", Uuid).First(&outlet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet not found")
		}
		log.Printf("Error finding outlet for update: %v", err)
		return nil, errors.New("failed to retrieve outlet for update")
	}

	// Update fields
	outlet.Name = updatedOutlet.Name
	outlet.Address = updatedOutlet.Address
	outlet.Type = updatedOutlet.Type

	if err := s.DB.Save(&outlet).Error; err != nil {
		log.Printf("Error updating outlet: %v", err)
		return nil, errors.New("failed to update outlet")
	}
	return &outlet, nil
}

func (s *OutletService) DeleteOutlet(Uuid uuid.UUID) error {
	if err := s.DB.Where("uuid = ?", Uuid).Delete(&models.Outlet{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("outlet not found")
		}
		log.Printf("Error deleting outlet: %v", err)
		return errors.New("failed to delete outlet")
	}
	return nil
}
