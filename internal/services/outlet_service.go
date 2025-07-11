package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type OutletService struct {
	DB                 *gorm.DB
	UserContextService *UserContextService
}

func NewOutletService(db *gorm.DB, userContextService *UserContextService) *OutletService {
	return &OutletService{DB: db, UserContextService: userContextService}
}

func (s *OutletService) GetAllOutlets(userID uint) ([]models.Outlet, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var outlets []models.Outlet
	if err := s.DB.Where("user_id = ?", ownerID).Find(&outlets).Error; err != nil {
		log.Printf("Error getting all outlets: %v", err)
		return nil, errors.New("failed to retrieve outlets")
	}
	return outlets, nil
}

func (s *OutletService) GetOutletByUuid(Uuid uuid.UUID, userID uint) (*dtos.OutletResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", Uuid, ownerID).First(&outlet).Error; err != nil {
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

func (s *OutletService) CreateOutlet(req *dtos.OutletCreateRequest, userID uint) (*dtos.OutletResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	outlet := &models.Outlet{
		Name:    req.Name,
		Address: req.Address,
		Type:    req.Type,
		UserID:  ownerID,
	}
	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(outlet).Error; err != nil {
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

func (s *OutletService) UpdateOutlet(Uuid uuid.UUID, req *dtos.OutletUpdateRequest, userID uint) (*dtos.OutletResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", Uuid, ownerID).First(&outlet).Error; err != nil {
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

func (s *OutletService) DeleteOutlet(Uuid uuid.UUID, userID uint) error {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return err
	}
	if err := s.DB.Where("uuid = ? AND user_id = ?", Uuid, ownerID).Delete(&models.Outlet{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("outlet not found")
		}
		log.Printf("Error deleting outlet: %v", err)
		return errors.New("failed to delete outlet")
	}
	return nil
}
