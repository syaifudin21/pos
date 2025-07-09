package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type StockService struct {
	DB                   *gorm.DB
	UserContextService   *UserContextService
	StockMovementService *StockMovementService
}

func NewStockService(db *gorm.DB, userContextService *UserContextService, stockMovementService *StockMovementService) *StockService {
	return &StockService{DB: db, UserContextService: userContextService, StockMovementService: stockMovementService}
}

// GetOutletStocks retrieves all stocks for a given outlet, including variants.
func (s *StockService) GetOutletStocks(outletUuid uuid.UUID, userID uint) ([]dtos.StockResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var stocks []models.Stock
	if err := s.DB.Preload("Product").Preload("ProductVariant.Product").Where("outlet_id = ? AND user_id = ?", outlet.ID, ownerID).Find(&stocks).Error; err != nil {
		log.Printf("Error getting outlet stocks: %v", err)
		return nil, errors.New("failed to retrieve outlet stocks")
	}

	var stockResponses []dtos.StockResponse
	for _, stock := range stocks {
		if stock.ProductID != nil && stock.Product != nil {
			stockResponses = append(stockResponses, dtos.StockResponse{
				ProductUuid: stock.Product.Uuid,
				ProductName: stock.Product.Name,
				ProductSku:  stock.Product.SKU,
				Quantity:    stock.Quantity,
			})
		} else if stock.ProductVariantID != nil && stock.ProductVariant != nil && stock.ProductVariant.Product.ID != 0 {
			stockResponses = append(stockResponses, dtos.StockResponse{
				ProductUuid:        stock.ProductVariant.Product.Uuid,
				ProductName:        stock.ProductVariant.Product.Name,
				ProductVariantUuid: &stock.ProductVariant.Uuid,
				VariantName:        stock.ProductVariant.Name,
				VariantSku:         stock.ProductVariant.SKU,
				Quantity:           stock.Quantity,
			})
		}
	}
	return stockResponses, nil
}

// UpdateStock handles updating stock for both simple products and variants.
func (s *StockService) UpdateStock(req dtos.UpdateStockRequest, outletUuid uuid.UUID, userID uint) (*dtos.StockResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var stock models.Stock
	var product *models.Product
	var variant *models.ProductVariant
	var query *gorm.DB

	if req.ProductVariantUuid != uuid.Nil {
		// Find variant
		if err := s.DB.Where("uuid = ? AND user_id = ?", req.ProductVariantUuid, ownerID).First(&variant).Error; err != nil {
			return nil, errors.New("product variant not found")
		}
		query = s.DB.Where("outlet_id = ? AND product_variant_id = ? AND user_id = ?", outlet.ID, variant.ID, ownerID)
	} else if req.ProductUuid != uuid.Nil {
		// Find product
		if err := s.DB.Where("uuid = ? AND user_id = ?", req.ProductUuid, ownerID).First(&product).Error; err != nil {
			return nil, errors.New("product not found")
		}
		query = s.DB.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outlet.ID, product.ID, ownerID)
	} else {
		return nil, errors.New("product_uuid or product_variant_uuid is required")
	}

	oldQuantity := 0.0
	quantityChange := req.Quantity

	err = query.First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new stock
			stock = models.Stock{
				OutletID: outlet.ID,
				Quantity: req.Quantity,
				UserID:   ownerID,
			}
			if variant != nil {
				stock.ProductVariantID = &variant.ID
			} else {
				stock.ProductID = &product.ID
			}
			if err := s.DB.Create(&stock).Error; err != nil {
				log.Printf("Error creating stock: %v", err)
				return nil, errors.New("failed to create stock")
			}
		} else {
			log.Printf("Error finding stock for update: %v", err)
			return nil, errors.New("failed to retrieve stock for update")
		}
	} else {
		// Update existing stock
		oldQuantity = stock.Quantity
		stock.Quantity = req.Quantity
		quantityChange = req.Quantity - oldQuantity
		if err := s.DB.Save(&stock).Error; err != nil {
			log.Printf("Error updating stock: %v", err)
			return nil, errors.New("failed to update stock")
		}
	}

	// Record stock movement
	if quantityChange != 0 {
		movement := &models.StockMovement{
			OutletID:       outlet.ID,
			QuantityChange: int(quantityChange),
			MovementType:   "Adjustment",
			Description:    stringPtr("Direct stock update"),
		}
		if variant != nil {
			movement.ProductVariantID = &variant.ID
		} else if product != nil {
			movement.ProductID = &product.ID
		}
		if err := s.StockMovementService.CreateStockMovement(movement); err != nil {
			log.Printf("Error recording stock movement: %v", err)
		}
	}

	// Build response
	resp := &dtos.StockResponse{Quantity: stock.Quantity}
	if variant != nil {
		// Preload parent product for name
		s.DB.Model(variant).Association("Product").Find(&variant.Product)
		resp.ProductUuid = variant.Product.Uuid
		resp.ProductName = variant.Product.Name
		resp.ProductVariantUuid = &variant.Uuid
		resp.VariantName = variant.Name
		resp.VariantSku = variant.SKU
	} else if product != nil {
		resp.ProductUuid = product.Uuid
		resp.ProductName = product.Name
		resp.ProductSku = product.SKU
	}

	return resp, nil
}

func (s *StockService) DeductStockForSale(tx *gorm.DB, outletID uint, productID *uint, productVariantID *uint, quantity float64, userID uint) error {
	var stock models.Stock
	var query *gorm.DB

	if productVariantID != nil {
		query = tx.Where("outlet_id = ? AND product_variant_id = ? AND user_id = ?", outletID, *productVariantID, userID)
	} else if productID != nil {
		query = tx.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outletID, *productID, userID)
	} else {
		return errors.New("product_id or product_variant_id is required for stock deduction")
	}

	if err := query.First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("stock not found")
		}
		return err
	}

	if stock.Quantity < quantity {
		return errors.New("insufficient stock")
	}

	stock.Quantity -= quantity
	if err := tx.Save(&stock).Error; err != nil {
		return err
	}

	// Record stock movement
	movement := &models.StockMovement{
		OutletID:       outletID,
		ProductID:      productID,
		ProductVariantID: productVariantID,
		QuantityChange: int(-quantity),
		MovementType:   "Order",
		Description:    stringPtr("Deduction for sale"),
	}
	return s.StockMovementService.CreateStockMovementWithTx(tx, movement)
}

// stringPtr is a helper function to return a pointer to a string.
func stringPtr(s string) *string {
	return &s
}
