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

// GetOwnerID retrieves the owner's ID for a given user.
// If the user is a manager or cashier, it returns their creator's ID.
// Otherwise, it returns the user's own ID.
func (s *OutletService) GetOwnerID(userID uint) (uint, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		log.Printf("Error finding user: %v", err)
		return 0, errors.New("user not found")
	}

	if (user.Role == "manager" || user.Role == "cashier") && user.CreatorID != nil {
		return *user.CreatorID, nil
	}

	return userID, nil
}

func (s *OutletService) GetAllOutlets(userID uint) ([]models.Outlet, error) {
	ownerID, err := s.GetOwnerID(userID)
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
	ownerID, err := s.GetOwnerID(userID)
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
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	outlet := &models.Outlet{
		Name:    req.Name,
		Address: req.Address,
		Type:    req.Type,
		UserID:  ownerID,
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

func (s *OutletService) UpdateOutlet(Uuid uuid.UUID, req *dtos.OutletUpdateRequest, userID uint) (*dtos.OutletResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
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
	ownerID, err := s.GetOwnerID(userID)
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
