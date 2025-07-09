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
	DB                 *gorm.DB
	UserContextService *UserContextService
}

func NewProductService(db *gorm.DB, userContextService *UserContextService) *ProductService {
	return &ProductService{DB: db, UserContextService: userContextService}
}

func (s *ProductService) GetAllProducts(userID uint, productType string) ([]models.Product, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var products []models.Product
	query := s.DB.Preload("Variants").Where("user_id = ?", ownerID)

	if productType != "" {
		query = query.Where("type = ?", productType)
	}

	if err := query.Find(&products).Error; err != nil {
		log.Printf("Error getting all products: %v", err)
		return nil, errors.New("failed to retrieve products")
	}
	return products, nil
}

func (s *ProductService) GetProductByUuid(Uuid uuid.UUID, userID uint) (*dtos.ProductResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var product models.Product
	if err := s.DB.Preload("Variants").Preload("Recipes.Component").Where("uuid = ? AND user_id = ?", Uuid, ownerID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		log.Printf("Error getting product by Uuid: %v", err)
		return nil, errors.New("failed to retrieve product")
	}

	variantResponses := []dtos.ProductVariantResponse{}
	for _, v := range product.Variants {
		variantResponses = append(variantResponses, dtos.ProductVariantResponse{
			ID:    v.ID,
			Uuid:  v.Uuid,
			Name:  v.Name,
			SKU:   v.SKU,
			Price: v.Price,
		})
	}

	recipeResponses := []dtos.RecipeResponse{}
	if product.Type == "fnb_main_product" {
		for _, r := range product.Recipes {
			if r.Component.ID != 0 { // Check if component is loaded
				recipeResponses = append(recipeResponses, dtos.RecipeResponse{
					Uuid:        r.Uuid,
					ComponentID: r.ComponentID,
					ComponentName:   r.Component.Name,
					Quantity:    r.Quantity,
				})
			}
		}
	}

	return &dtos.ProductResponse{
		ID:          product.ID,
		Uuid:        product.Uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Type:        product.Type,
		Variants:    variantResponses,
		Recipes:     recipeResponses,
	}, nil
}

func (s *ProductService) CreateProduct(req *dtos.ProductCreateRequest, userID uint) (*dtos.ProductResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
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

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(product).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating product: %v", err)
		return nil, errors.New("failed to create product")
	}

	var variantResponses []dtos.ProductVariantResponse
	if len(req.Variants) > 0 {
		for _, v := range req.Variants {
			productVariant := &models.ProductVariant{
				ProductID: product.ID,
				Name:      v.Name,
				SKU:       v.SKU,
				Price:     v.Price,
				UserID:    ownerID,
			}
			if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(productVariant).Error; err != nil {
				tx.Rollback()
				log.Printf("Error creating product variant: %v", err)
				return nil, errors.New("failed to create product variant")
			}
			variantResponses = append(variantResponses, dtos.ProductVariantResponse{
				ID:    productVariant.ID,
				Uuid:  productVariant.Uuid,
				Name:  productVariant.Name,
				SKU:   productVariant.SKU,
				Price: productVariant.Price,
			})
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, errors.New("failed to create product with variants")
	}

	return &dtos.ProductResponse{
		ID:          product.ID,
		Uuid:        product.Uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Type:        product.Type,
		Variants:    variantResponses,
	}, nil
}

func (s *ProductService) UpdateProduct(Uuid uuid.UUID, req *dtos.ProductUpdateRequest, userID uint) (*dtos.ProductResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var product models.Product
	if err := tx.Where("uuid = ? AND user_id = ?", Uuid, ownerID).First(&product).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		log.Printf("Error finding product for update: %v", err)
		return nil, errors.New("failed to retrieve product for update")
	}

	// Update main product fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.SKU = req.SKU
	product.Type = req.Type

	if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&product).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating product: %v", err)
		return nil, errors.New("failed to update product")
	}

	//- Hapus varian lama
	if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductVariant{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting old variants: %v", err)
		return nil, errors.New("failed to update variants")
	}

	// Buat varian baru
	var variantResponses []dtos.ProductVariantResponse
	if len(req.Variants) > 0 {
		for _, v := range req.Variants {
			productVariant := &models.ProductVariant{
				ProductID: product.ID,
				Name:      v.Name,
				SKU:       v.SKU,
				Price:     v.Price,
				UserID:    ownerID,
			}
			if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(productVariant).Error; err != nil {
				tx.Rollback()
				log.Printf("Error creating new variant: %v", err)
				return nil, errors.New("failed to create new variant")
			}
			variantResponses = append(variantResponses, dtos.ProductVariantResponse{
				ID:    productVariant.ID,
				Uuid:  productVariant.Uuid,
				Name:  productVariant.Name,
				SKU:   productVariant.SKU,
				Price: productVariant.Price,
			})
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for update: %v", err)
		return nil, errors.New("failed to update product with variants")
	}

	return &dtos.ProductResponse{
		ID:          product.ID,
		Uuid:        product.Uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Type:        product.Type,
		Variants:    variantResponses,
	}, nil
}

func (s *ProductService) DeleteProduct(Uuid uuid.UUID, userID uint) error {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
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
	ownerID, err := s.UserContextService.GetOwnerID(userID)
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
