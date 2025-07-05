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

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

// GetOwnerID retrieves the owner's ID for a given user.
// If the user is a manager or cashier, it returns their creator's ID.
// Otherwise, it returns the user's own ID.
func (s *ProductService) GetOwnerID(userID uint) (uint, error) {
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

func (s *ProductService) GetAllProducts(userID uint) ([]models.Product, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var products []models.Product
	if err := s.DB.Where("user_id = ?", ownerID).Find(&products).Error; err != nil {
		log.Printf("Error getting all products: %v", err)
		return nil, errors.New("failed to retrieve products")
	}
	return products, nil
}

func (s *ProductService) GetProductByUuid(Uuid uuid.UUID, userID uint) (*dtos.ProductResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var product models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", Uuid, ownerID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		log.Printf("Error getting product by Uuid: %v", err)
		return nil, errors.New("failed to retrieve product")
	}
	return &dtos.ProductResponse{
		ID:          product.ID,
		Uuid:        product.Uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Type:        product.Type,
	}, nil
}

func (s *ProductService) CreateProduct(req *dtos.ProductCreateRequest, userID uint) (*dtos.ProductResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
		Type:        req.Type,
		UserID:      ownerID,
	}
	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(product).Error; err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, errors.New("failed to create product")
	}
	return &dtos.ProductResponse{
		ID:          product.ID,
		Uuid:        product.Uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Type:        product.Type,
	}, nil
}

func (s *ProductService) UpdateProduct(Uuid uuid.UUID, req *dtos.ProductUpdateRequest, userID uint) (*dtos.ProductResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var product models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", Uuid, ownerID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		log.Printf("Error finding product for update: %v", err)
		return nil, errors.New("failed to retrieve product for update")
	}

	// Update fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.SKU = req.SKU
	product.Type = req.Type

	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&product).Error; err != nil {
		log.Printf("Error updating product: %v", err)
		return nil, errors.New("failed to update product")
	}
	return &dtos.ProductResponse{
		ID:          product.ID,
		Uuid:        product.Uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Type:        product.Type,
	}, nil
}

func (s *ProductService) DeleteProduct(Uuid uuid.UUID, userID uint) error {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return err
	}
	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Where("uuid = ? AND user_id = ?", Uuid, ownerID).Delete(&models.Product{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		log.Printf("Error deleting product: %v", err)
		return errors.New("failed to delete product")
	}
	return nil
}

// GetProductsByOutlet retrieves all products available in a specific outlet (i.e., have stock).
func (s *ProductService) GetProductsByOutlet(outletUuid uuid.UUID, userID uint) ([]dtos.ProductOutletResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var products []dtos.ProductOutletResponse

	// Find the outlet first
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	// Join products with stocks to get products available in the outlet
	if err := s.DB.Table("products").
		Select("products.uuid as product_uuid, products.name as product_name, products.sku as product_sku, products.price, products.type, stocks.quantity").
		Joins("JOIN stocks ON products.id = stocks.product_id").
		Where("stocks.outlet_id = ? AND stocks.quantity > 0 AND products.user_id = ?", outlet.ID, ownerID).
		Find(&products).Error; err != nil {
		log.Printf("Error getting products by outlet: %v", err)
		return nil, errors.New("failed to retrieve products for outlet")
	}

	return products, nil
}
