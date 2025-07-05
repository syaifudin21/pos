package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
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

func (s *OutletService) GetOutletByUuid(Uuid uuid.UUID) (*dtos.OutletResponse, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", Uuid).First(&outlet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet not found")
		}
		log.Printf("Error getting outlet by Uuid: %v", err)
		return nil, errors.New("failed to retrieve outlet")
	}
	return &dtos.OutletResponse{
		ID:      outlet.ID,
		Uuid:    outlet.Uuid,
		Name:    outlet.Name,
		Address: outlet.Address,
		Type:    outlet.Type,
	}, nil
}

func (s *OutletService) CreateOutlet(req *dtos.OutletCreateRequest) (*dtos.OutletResponse, error) {
	outlet := &models.Outlet{
		Name:    req.Name,
		Address: req.Address,
		Type:    req.Type,
	}
	if err := s.DB.Create(outlet).Error; err != nil {
		log.Printf("Error creating outlet: %v", err)
		return nil, errors.New("failed to create outlet")
	}
	return &dtos.OutletResponse{
		ID:      outlet.ID,
		Uuid:    outlet.Uuid,
		Name:    outlet.Name,
		Address: outlet.Address,
		Type:    outlet.Type,
	}, nil
}

func (s *OutletService) UpdateOutlet(Uuid uuid.UUID, req *dtos.OutletUpdateRequest) (*dtos.OutletResponse, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", Uuid).First(&outlet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet not found")
		}
		log.Printf("Error finding outlet for update: %v", err)
		return nil, errors.New("failed to retrieve outlet for update")
	}

	// Update fields
	outlet.Name = req.Name
	outlet.Address = req.Address
	outlet.Type = req.Type

	if err := s.DB.Save(&outlet).Error; err != nil {
		log.Printf("Error updating outlet: %v", err)
		return nil, errors.New("failed to update outlet")
	}
	return &dtos.OutletResponse{
		ID:      outlet.ID,
		Uuid:    outlet.Uuid,
		Name:    outlet.Name,
		Address: outlet.Address,
		Type:    outlet.Type,
	}, nil
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
