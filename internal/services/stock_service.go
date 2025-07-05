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

type StockService struct {
	DB *gorm.DB
}

func NewStockService(db *gorm.DB) *StockService {
	return &StockService{DB: db}
}

// GetOwnerID retrieves the owner's ID for a given user.
// If the user is a manager or cashier, it returns their creator's ID.
// Otherwise, it returns the user's own ID.
func (s *StockService) GetOwnerID(userID uint) (uint, error) {
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

// GetStockByOutletAndProduct retrieves stock for a specific product in an outlet.
func (s *StockService) GetStockByOutletAndProduct(outletUuid, productUuid uuid.UUID, userID uint) (*dtos.StockResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var stock models.Stock
	var outlet models.Outlet
	var product models.Product

	err = s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error
	if err != nil {
		return nil, errors.New("outlet not found")
	}

	err = s.DB.Where("uuid = ? AND user_id = ?", productUuid, ownerID).First(&product).Error
	if err != nil {
		return nil, errors.New("product not found")
	}

	err = s.DB.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outlet.ID, product.ID, ownerID).First(&stock).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stock not found for this outlet and product")
		}
		log.Printf("Error getting stock by outlet and product: %v", err)
		return nil, errors.New("failed to retrieve stock")
	}
	return &dtos.StockResponse{
		ProductUuid: product.Uuid,
		ProductName: product.Name,
		ProductSku:  product.SKU,
		Quantity:    stock.Quantity,
	}, nil
}

// GetOutletStocks retrieves all stocks for a given outlet.
func (s *StockService) GetOutletStocks(outletUuid uuid.UUID, userID uint) ([]dtos.StockResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var stocks []models.Stock
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	if err := s.DB.Preload("Product").Where("outlet_id = ? AND user_id = ?", outlet.ID, ownerID).Find(&stocks).Error; err != nil {
		log.Printf("Error getting outlet stocks: %v", err)
		return nil, errors.New("failed to retrieve outlet stocks")
	}

	var stockResponses []dtos.StockResponse
	for _, stock := range stocks {
		stockResponses = append(stockResponses, dtos.StockResponse{
			ProductUuid: stock.Product.Uuid,
			ProductName: stock.Product.Name,
			ProductSku:  stock.Product.SKU,
			Quantity:    stock.Quantity,
		})
	}
	return stockResponses, nil
}

// UpdateStock updates the quantity of a product in an outlet.
// This is a direct update, useful for initial setup or corrections.
func (s *StockService) UpdateStock(outletUuid, productUuid uuid.UUID, quantity float64, userID uint) (*dtos.StockResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var product models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", productUuid, ownerID).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}

	var stock models.Stock
	if err := s.DB.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outlet.ID, product.ID, ownerID).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new stock entry if not found
			stock = models.Stock{
				OutletID:  outlet.ID,
				ProductID: product.ID,
				Quantity:  quantity,
				UserID:    ownerID,
			}
			if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&stock).Error; err != nil {
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
		if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&stock).Error; err != nil {
			log.Printf("Error updating stock: %v", err)
			return nil, errors.New("failed to update stock")
		}
	}

	return &dtos.StockResponse{
		ProductUuid: product.Uuid,
		ProductName: product.Name,
		ProductSku:  product.SKU,
		Quantity:    stock.Quantity,
	}, nil
}

// DeductStockForSale handles stock deduction based on product type.
// For FnB main products, it deducts from components based on recipe.
func (s *StockService) DeductStockForSale(outletUuid, productUuid uuid.UUID, quantity float64, userID uint) error {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return err
	}
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return errors.New("outlet not found")
	}

	var product models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", productUuid, ownerID).First(&product).Error; err != nil {
		return errors.New("product not found")
	}

	log.Printf("DeductStockForSale: Processing product %s (Type: %s) for outlet %s, quantity: %f", product.Name, product.Type, outlet.Name, quantity)

	if product.Type == "fnb_main_product" {
		// Deduct components based on recipe
		var recipes []models.Recipe
		if err := s.DB.Where("main_product_id = ? AND user_id = ?", product.ID, ownerID).Find(&recipes).Error; err != nil {
			log.Printf("Error finding recipes for product %s: %v", product.Name, err)
			return errors.New("failed to retrieve product recipe")
		}

		if len(recipes) == 0 {
			return errors.New("FnB main product has no defined recipe")
		}

		for _, recipe := range recipes {
			requiredComponentQuantity := recipe.Quantity * quantity
			var componentStock models.Stock
			if err := s.DB.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outlet.ID, recipe.ComponentID, ownerID).First(&componentStock).Error; err != nil {
				log.Printf("DeductStockForSale: Component stock not found for product %s (component of %s) in outlet %s. Error: %v", recipe.Component.Name, product.Name, outlet.Name, err)
				return errors.New("component stock not found")
			}

			log.Printf("DeductStockForSale: Before deduction - Product %s (component of %s), Current Stock: %f", recipe.Component.Name, product.Name, componentStock.Quantity)

			if componentStock.Quantity < requiredComponentQuantity {
				log.Printf("DeductStockForSale: Insufficient stock for component %s. Available: %f, Required: %f", recipe.Component.Name, componentStock.Quantity, requiredComponentQuantity)
				return errors.New("insufficient stock for components")
			}

			componentStock.Quantity -= requiredComponentQuantity
			if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&componentStock).Error; err != nil {
				log.Printf("Error deducting component stock for product %s: %v", recipe.Component.Name, err)
				return errors.New("failed to deduct component stock")
			}
			log.Printf("DeductStockForSale: After deduction - Product %s, New Stock: %f", recipe.Component.Name, product.Name, componentStock.Quantity)
		}
	} else if product.Type == "retail_item" {
		var retailStock models.Stock
		if err := s.DB.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outlet.ID, product.ID, userID).First(&retailStock).Error; err != nil {
			log.Printf("DeductStockForSale: Retail item stock not found for product %s in outlet %s. Error: %v", product.Name, outlet.Name, err)
			return errors.New("retail item stock not found")
		}

		log.Printf("DeductStockForSale: Before deduction - Product %s, Current Stock: %f", product.Name, retailStock.Quantity)

		if retailStock.Quantity < quantity {
			log.Printf("DeductStockForSale: Insufficient stock for retail item %s. Available: %f, Required: %f", product.Name, retailStock.Quantity, quantity)
			return errors.New("insufficient stock for retail item")
		}

		retailStock.Quantity -= quantity
		if err := s.DB.Save(&retailStock).Error; err != nil {
			log.Printf("Error deducting retail item stock for product %s: %v", product.Name, err)
			return errors.New("failed to deduct retail item stock")
		}
		log.Printf("DeductStockForSale: After deduction - Product %s, New Stock: %f", product.Name, retailStock.Quantity)
	}

	return nil
}

func (s *StockService) UpdateGlobalStock(outletUuid, productUuid uuid.UUID, quantity float64, userID uint) (*dtos.StockResponse, error) {
	return s.UpdateStock(outletUuid, productUuid, quantity, userID)
}

// AdjustStock adds or subtracts quantity from an existing stock entry.
// If stock does not exist, it creates a new one.
func (s *StockService) AdjustStock(outletUuid, productUuid uuid.UUID, quantityChange float64, userID uint) (*dtos.StockResponse, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, userID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var product models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", productUuid, userID).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}

	var stock models.Stock
	if err := s.DB.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outlet.ID, product.ID, userID).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new stock entry if not found
			stock = models.Stock{
				OutletID:  outlet.ID,
				ProductID: product.ID,
				Quantity:  quantityChange,
				UserID:    userID,
			}
			if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&stock).Error; err != nil {
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

	return &dtos.StockResponse{
		ProductUuid: product.Uuid,
		ProductName: product.Name,
		ProductSku:  product.SKU,
		Quantity:    stock.Quantity,
	}, nil
}
