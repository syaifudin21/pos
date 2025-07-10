package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type ProductAddOnService struct {
	DB                 *gorm.DB
	UserContextService *UserContextService
}

func NewProductAddOnService(db *gorm.DB, userContextService *UserContextService) *ProductAddOnService {
	return &ProductAddOnService{DB: db, UserContextService: userContextService}
}

func (s *ProductAddOnService) CreateProductAddOn(req *dtos.ProductAddOnRequest, userID uint) (*dtos.ProductAddOnResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", req.ProductID, ownerID).First(&product).Error; err != nil {
		return nil, errors.New("main product not found")
	}

	var addOnProduct models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ? AND type = ?", req.AddOnID, ownerID, "add_on").First(&addOnProduct).Error; err != nil {
		return nil, errors.New("add-on product not found or not of type add_on")
	}

	// Check if add-on already exists for this product
	var existingAddOn models.ProductAddOn
	if err := s.DB.Where("product_id = ? AND add_on_id = ? AND user_id = ?", product.ID, addOnProduct.ID, ownerID).First(&existingAddOn).Error; err == nil {
		return nil, errors.New("add-on already exists for this product")
	}

	productAddOn := &models.ProductAddOn{
		ProductID:   product.ID,
		AddOnID:     addOnProduct.ID,
		Price:       req.Price,
		IsAvailable: true,
		UserID:      ownerID,
	}

	if err := s.DB.Create(productAddOn).Error; err != nil {
		log.Printf("Error creating product add-on: %v", err)
		return nil, errors.New("failed to create product add-on")
	}

	return &dtos.ProductAddOnResponse{
		Uuid:        productAddOn.Uuid,
		AddOnName:   addOnProduct.Name,
		Price:       productAddOn.Price,
		IsAvailable: productAddOn.IsAvailable,
	}, nil
}

func (s *ProductAddOnService) GetProductAddOnsByProductID(productUuid uuid.UUID, userID uint) ([]dtos.ProductAddOnResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", productUuid, ownerID).First(&product).Error; err != nil {
		return nil, errors.New("main product not found")
	}

	var productAddOns []models.ProductAddOn
	if err := s.DB.Preload("AddOn").Where("product_id = ? AND user_id = ?", product.ID, ownerID).Find(&productAddOns).Error; err != nil {
		log.Printf("Error getting product add-ons: %v", err)
		return nil, errors.New("failed to retrieve product add-ons")
	}

	var responses []dtos.ProductAddOnResponse
	for _, pao := range productAddOns {
		responses = append(responses, dtos.ProductAddOnResponse{
			Uuid:        pao.Uuid,
			AddOnName:   pao.AddOn.Name,
			Price:       pao.Price,
			IsAvailable: pao.IsAvailable,
		})
	}
	return responses, nil
}

func (s *ProductAddOnService) DeleteProductAddOn(productAddOnUuid uuid.UUID, userID uint) error {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return err
	}

	var productAddOn models.ProductAddOn
	if err := s.DB.Where("uuid = ? AND user_id = ?", productAddOnUuid, ownerID).First(&productAddOn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product add-on not found")
		}
		log.Printf("Error finding product add-on for deletion: %v", err)
		return errors.New("failed to retrieve product add-on for deletion")
	}

	if err := s.DB.Delete(&productAddOn).Error; err != nil {
		log.Printf("Error deleting product add-on: %v", err)
		return errors.New("failed to delete product add-on")
	}
	return nil
}
