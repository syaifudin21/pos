package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type StockService struct {
	DB *gorm.DB
}

func NewStockService(db *gorm.DB) *StockService {
	return &StockService{DB: db}
}

// GetStockByOutletAndProduct retrieves stock for a specific product in an outlet.
func (s *StockService) GetStockByOutletAndProduct(outletUuid, productUuid uuid.UUID) (*models.Stock, error) {
	var stock models.Stock
	err := s.DB.Preload("Outlet").Preload("Product").
		Joins("JOIN outlets ON stocks.outlet_id = outlets.id").
		Joins("JOIN products ON stocks.product_id = products.id").
		Where("outlets.uuid = ? AND products.uuid = ?", outletUuid, productUuid).
		First(&stock).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stock not found for this outlet and product")
		}
		log.Printf("Error getting stock by outlet and product: %v", err)
		return nil, errors.New("failed to retrieve stock")
	}
	return &stock, nil
}

// GetOutletStocks retrieves all stocks for a given outlet.
func (s *StockService) GetOutletStocks(outletUuid uuid.UUID) ([]models.Stock, error) {
	var stocks []models.Stock
	err := s.DB.Preload("Product").
		Joins("JOIN outlets ON stocks.outlet_id = outlets.id").
		Where("outlets.uuid = ?", outletUuid).
		Find(&stocks).Error

	if err != nil {
		log.Printf("Error getting outlet stocks: %v", err)
		return nil, errors.New("failed to retrieve outlet stocks")
	}
	return stocks, nil
}

// UpdateStock updates the quantity of a product in an outlet.
// This is a direct update, useful for initial setup or corrections.
func (s *StockService) UpdateStock(outletUuid, productUuid uuid.UUID, quantity float64) (*models.Stock, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", outletUuid).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var product models.Product
	if err := s.DB.Where("uuid = ?", productUuid).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}

	var stock models.Stock
	if err := s.DB.Where("outlet_id = ? AND product_id = ?", outlet.ID, product.ID).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new stock entry if not found
			stock = models.Stock{
				OutletID:  outlet.ID,
				ProductID: product.ID,
				Quantity:  quantity,
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
		stock.Quantity = quantity
		if err := s.DB.Save(&stock).Error; err != nil {
			log.Printf("Error updating stock: %v", err)
			return nil, errors.New("failed to update stock")
		}
	}

	// Reload stock with associations for response
	s.DB.Preload("Outlet").Preload("Product").First(&stock, stock.ID)
	return &stock, nil
}

// DeductStockForSale handles stock deduction based on product type.
// For FnB main products, it deducts from components based on recipe.
func (s *StockService) DeductStockForSale(OutletUuid, productuuid uuid.UUID, quantity float64) error {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", OutletUuid).First(&outlet).Error; err != nil {
		return errors.New("outlet not found")
	}

	var product models.Product
	if err := s.DB.Where("uuid = ?", productuuid).First(&product).Error; err != nil {
		return errors.New("product not found")
	}

	if product.Type == "fnb_main_product" {
		// Deduct components based on recipe
		var recipes []models.Recipe
		if err := s.DB.Where("main_product_id = ?", product.ID).Find(&recipes).Error; err != nil {
			log.Printf("Error finding recipes for product %s: %v", product.Name, err)
			return errors.New("failed to retrieve product recipe")
		}

		if len(recipes) == 0 {
			return errors.New("FnB main product has no defined recipe")
		}

		for _, recipe := range recipes {
			requiredComponentQuantity := recipe.Quantity * quantity
			var componentStock models.Stock
			if err := s.DB.Where("outlet_id = ? AND product_id = ?", outlet.ID, recipe.ComponentID).First(&componentStock).Error; err != nil {
				return errors.New("component stock not found")
			}

			if componentStock.Quantity < requiredComponentQuantity {
				return errors.New("insufficient stock for components")
			}

			componentStock.Quantity -= requiredComponentQuantity
			if err := s.DB.Save(&componentStock).Error; err != nil {
				log.Printf("Error deducting component stock: %v", err)
				return errors.New("failed to deduct component stock")
			}
		}
	} else {
		// Deduct directly for retail items or FnB components
		var stock models.Stock
		if err := s.DB.Where("outlet_id = ? AND product_id = ?", outlet.ID, product.ID).First(&stock).Error; err != nil {
			return errors.New("stock not found")
		}

		if stock.Quantity < quantity {
			return errors.New("insufficient stock")
		}

		stock.Quantity -= quantity
		if err := s.DB.Save(&stock).Error; err != nil {
			log.Printf("Error deducting stock: %v", err)
			return errors.New("failed to deduct stock")
		}
	}

	return nil
}

func (s *StockService) UpdateGlobalStock(outletUuid, productUuid uuid.UUID, quantity float64) (*models.Stock, error) {
	return s.UpdateStock(outletUuid, productUuid, quantity)
}

// AdjustStock adds or subtracts quantity from an existing stock entry.
// If stock does not exist, it creates a new one.
func (s *StockService) AdjustStock(outletUuid, productUuid uuid.UUID, quantityChange float64) (*models.Stock, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ?", outletUuid).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var product models.Product
	if err := s.DB.Where("uuid = ?", productUuid).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}

	var stock models.Stock
	if err := s.DB.Where("outlet_id = ? AND product_id = ?", outlet.ID, product.ID).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new stock entry if not found
			stock = models.Stock{
				OutletID:  outlet.ID,
				ProductID: product.ID,
				Quantity:  quantityChange,
			}
			if err := s.DB.Create(&stock).Error; err != nil {
				log.Printf("Error creating stock: %v", err)
				return nil, errors.New("failed to create stock")
			}
		} else {
			log.Printf("Error finding stock for adjustment: %v", err)
			return nil, errors.New("failed to retrieve stock for adjustment")
		}
	} else {
		// Adjust existing stock
		stock.Quantity += quantityChange
		if err := s.DB.Save(&stock).Error; err != nil {
			log.Printf("Error adjusting stock: %v", err)
			return nil, errors.New("failed to adjust stock")
		}
	}

	// Reload stock with associations for response
	s.DB.Preload("Outlet").Preload("Product").First(&stock, stock.ID)
	return &stock, nil
}
