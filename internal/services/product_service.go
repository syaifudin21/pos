package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := s.DB.Find(&products).Error; err != nil {
		log.Printf("Error getting all products: %v", err)
		return nil, errors.New("failed to retrieve products")
	}
	return products, nil
}

func (s *ProductService) GetProductByUuid(Uuid uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := s.DB.Where("uuid = ?", Uuid).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		log.Printf("Error getting product by Uuid: %v", err)
		return nil, errors.New("failed to retrieve product")
	}
	return &product, nil
}

func (s *ProductService) CreateProduct(product *models.Product) (*models.Product, error) {
	if err := s.DB.Create(product).Error; err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, errors.New("failed to create product")
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(Uuid uuid.UUID, updatedProduct *models.Product) (*models.Product, error) {
	var product models.Product
	if err := s.DB.Where("uuid = ?", Uuid).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		log.Printf("Error finding product for update: %v", err)
		return nil, errors.New("failed to retrieve product for update")
	}

	// Update fields
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	product.SKU = updatedProduct.SKU
	product.Type = updatedProduct.Type

	if err := s.DB.Save(&product).Error; err != nil {
		log.Printf("Error updating product: %v", err)
		return nil, errors.New("failed to update product")
	}
	return &product, nil
}

func (s *ProductService) DeleteProduct(Uuid uuid.UUID) error {
	if err := s.DB.Where("uuid = ?", Uuid).Delete(&models.Product{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		log.Printf("Error deleting product: %v", err)
		return errors.New("failed to delete product")
	}
	return nil
}
