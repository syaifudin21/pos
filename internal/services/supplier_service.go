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

type SupplierService struct {
	DB *gorm.DB
	UserContextService *UserContextService
}

func NewSupplierService(db *gorm.DB, userContextService *UserContextService) *SupplierService {
	return &SupplierService{DB: db, UserContextService: userContextService}
}

func (s *SupplierService) GetAllSuppliers(userID uint) ([]models.Supplier, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var suppliers []models.Supplier
	if err := s.DB.Where("user_id = ?", ownerID).Find(&suppliers).Error; err != nil {
		log.Printf("Error getting all suppliers: %v", err)
		return nil, errors.New("failed to retrieve suppliers")
	}
	return suppliers, nil
}

func (s *SupplierService) GetSupplierByuuid(uuid uuid.UUID, userID uint) (*dtos.SupplierResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var supplier models.Supplier
	if err := s.DB.Where("uuid = ? AND user_id = ?", uuid, ownerID).First(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("supplier not found")
		}
		log.Printf("Error getting supplier by uuid: %v", err)
		return nil, errors.New("failed to retrieve supplier")
	}
	return &dtos.SupplierResponse{
		ID:      supplier.ID,
		Uuid:    supplier.Uuid,
		Name:    supplier.Name,
		Contact: supplier.Contact,
		Address: supplier.Address,
	}, nil
}

func (s *SupplierService) CreateSupplier(req *dtos.CreateSupplierRequest, userID uint) (*dtos.SupplierResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	supplier := &models.Supplier{
		Name:    req.Name,
		Contact: req.Contact,
		Address: req.Address,
		UserID:  ownerID,
	}
	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(supplier).Error; err != nil {
		log.Printf("Error creating supplier: %v", err)
		return nil, errors.New("failed to create supplier")
	}
	return &dtos.SupplierResponse{
		ID:      supplier.ID,
		Uuid:    supplier.Uuid,
		Name:    supplier.Name,
		Contact: supplier.Contact,
		Address: supplier.Address,
	}, nil
}

func (s *SupplierService) UpdateSupplier(uuid uuid.UUID, req *dtos.UpdateSupplierRequest, userID uint) (*dtos.SupplierResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var supplier models.Supplier
	if err := s.DB.Where("uuid = ? AND user_id = ?", uuid, ownerID).First(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("supplier not found")
		}
		log.Printf("Error finding supplier for update: %v", err)
		return nil, errors.New("failed to retrieve supplier for update")
	}

	// Update fields
	supplier.Name = req.Name
	supplier.Contact = req.Contact
	supplier.Address = req.Address

	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&supplier).Error; err != nil {
		log.Printf("Error updating supplier: %v", err)
		return nil, errors.New("failed to update supplier")
	}
	return &dtos.SupplierResponse{
		ID:      supplier.ID,
		Uuid:    supplier.Uuid,
		Name:    supplier.Name,
		Contact: supplier.Contact,
		Address: supplier.Address,
	}, nil
}

func (s *SupplierService) DeleteSupplier(uuid uuid.UUID, userID uint) error {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return err
	}
	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Where("uuid = ? AND user_id = ?", uuid, ownerID).Delete(&models.Supplier{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("supplier not found")
		}
		log.Printf("Error deleting supplier: %v", err)
		return errors.New("failed to delete supplier")
	}
	return nil
}
