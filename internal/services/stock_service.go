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
func (s *StockService) GetOutletStocks(outletUuid uuid.UUID, userID uint, productType string, isForSale bool) ([]dtos.StockDetailResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var stocks []models.Stock
	query := s.DB.Where("stocks.outlet_id = ? AND stocks.user_id = ?", outlet.ID, ownerID)

	// Always include joins if either productType or isForSale is active
	if productType != "" || isForSale {
		query = query.Joins("LEFT JOIN products AS p_direct ON stocks.product_id = p_direct.id")
		query = query.Joins("LEFT JOIN product_variants AS pv ON stocks.product_variant_id = pv.id")
		query = query.Joins("LEFT JOIN products AS p_variant ON pv.product_id = p_variant.id")
	}

	if productType != "" {
		query = query.Where("(stocks.product_id IS NOT NULL AND p_direct.type = ?) OR (stocks.product_variant_id IS NOT NULL AND p_variant.type = ?)", productType, productType)
	}

	if isForSale {
		query = query.Where("(stocks.product_id IS NOT NULL AND (p_direct.type = ? OR p_direct.type = ?)) OR (stocks.product_variant_id IS NOT NULL AND (p_variant.type = ? OR p_variant.type = ?))", "retail_item", "fnb_main_product", "retail_item", "fnb_main_product")
	}

	// Preload the relationships after applying the joins and filters
	query = query.Preload("Product.Recipes.Component").Preload("Product.AddOns.AddOn").Preload("ProductVariant.Product.Recipes.Component").Preload("ProductVariant.Product.AddOns.AddOn")

	if err := query.Find(&stocks).Error; err != nil {
		log.Printf("Error getting outlet stocks: %v", err)
		return nil, errors.New("failed to retrieve outlet stocks")
	}

	var stockDetailResponses []dtos.StockDetailResponse
	for _, stock := range stocks {
		var productDetail dtos.ProductDetailResponse
		var recipes []dtos.RecipeResponse
		var addOns []dtos.ProductAddOnResponse
		var variants []dtos.ProductVariantResponse

		if stock.ProductID != nil && stock.Product != nil {
			productDetail.Uuid = stock.Product.Uuid
			productDetail.Name = stock.Product.Name
			productDetail.Description = stock.Product.Description
			productDetail.Price = stock.Product.Price
			productDetail.SKU = stock.Product.SKU
			productDetail.Type = stock.Product.Type

			// Map recipes
			if stock.Product.Type == "fnb_main_product" {
				for _, r := range stock.Product.Recipes {
					if r.Component.ID != 0 {
						recipes = append(recipes, dtos.RecipeResponse{
							Uuid:          r.Uuid,
							ComponentName: r.Component.Name,
							Quantity:      r.Quantity,
						})
					}
				}
			}
			productDetail.Recipes = recipes

			// Map add-ons
			for _, pao := range stock.Product.AddOns {
				if pao.AddOn.ID != 0 {
					addOns = append(addOns, dtos.ProductAddOnResponse{
						Uuid:        pao.Uuid,
						AddOnName:   pao.AddOn.Name,
						Price:       pao.Price,
						IsAvailable: pao.IsAvailable,
					})
				}
			}
			productDetail.AddOns = addOns

			// Map variants (if any)
			for _, v := range stock.Product.Variants {
				variants = append(variants, dtos.ProductVariantResponse{
					ID:    v.ID,
					Uuid:  v.Uuid,
					Name:  v.Name,
					SKU:   v.SKU,
					Price: v.Price,
				})
			}
			productDetail.Variants = variants

			stockDetailResponses = append(stockDetailResponses, dtos.StockDetailResponse{
				Uuid:        productDetail.Uuid,
				Name:        productDetail.Name,
				Description: productDetail.Description,
				Price:       productDetail.Price,
				SKU:         productDetail.SKU,
				Type:        productDetail.Type,
				Recipes:     productDetail.Recipes,
				AddOns:      productDetail.AddOns,
				Variants:    productDetail.Variants,
				Quantity:    stock.Quantity,
			})
		} else if stock.ProductVariantID != nil && stock.ProductVariant != nil && stock.ProductVariant.Product.ID != 0 {
			// For product variants, the main product details come from stock.ProductVariant.Product
			productDetail.Uuid = stock.ProductVariant.Product.Uuid
			productDetail.Name = stock.ProductVariant.Product.Name
			productDetail.Description = stock.ProductVariant.Product.Description
			productDetail.Price = stock.ProductVariant.Product.Price
			productDetail.SKU = stock.ProductVariant.Product.SKU
			productDetail.Type = stock.ProductVariant.Product.Type

			// Map recipes for the main product of the variant
			if stock.ProductVariant.Product.Type == "fnb_main_product" {
				for _, r := range stock.ProductVariant.Product.Recipes {
					if r.Component.ID != 0 {
						recipes = append(recipes, dtos.RecipeResponse{
							Uuid:          r.Uuid,
							ComponentName: r.Component.Name,
							Quantity:      r.Quantity,
						})
					}
				}
			}
			productDetail.Recipes = recipes

			// Map add-ons for the main product of the variant
			for _, pao := range stock.ProductVariant.Product.AddOns {
				if pao.AddOn.ID != 0 {
					addOns = append(addOns, dtos.ProductAddOnResponse{
						Uuid:        pao.Uuid,
						AddOnName:   pao.AddOn.Name,
						Price:       pao.Price,
						IsAvailable: pao.IsAvailable,
					})
				}
			}
			productDetail.AddOns = addOns

			// Map variants (only the current variant for this stock entry)
			variants = append(variants, dtos.ProductVariantResponse{
				ID:    stock.ProductVariant.ID,
				Uuid:  stock.ProductVariant.Uuid,
				Name:  stock.ProductVariant.Name,
				SKU:   stock.ProductVariant.SKU,
				Price: stock.ProductVariant.Price,
			})
			productDetail.Variants = variants

			stockDetailResponses = append(stockDetailResponses, dtos.StockDetailResponse{
				Uuid:        productDetail.Uuid,
				Name:        productDetail.Name,
				Description: productDetail.Description,
				Price:       productDetail.Price,
				SKU:         productDetail.SKU,
				Type:        productDetail.Type,
				Recipes:     productDetail.Recipes,
				AddOns:      productDetail.AddOns,
				Variants:    productDetail.Variants,
				Quantity:    stock.Quantity,
			})
		}
	}
	return stockDetailResponses, nil
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

// ProduceFNBProduct handles the production of F&B main products, deducting component stock.
func (s *StockService) ProduceFNBProduct(req dtos.FNBProductionRequest, outletUuid uuid.UUID, userID uint) (*dtos.FNBProductionResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var mainProduct models.Product
	if err := s.DB.Preload("Recipes.Component").Where("uuid = ? AND user_id = ? AND type = ?", req.FNBMainProductUuid, ownerID, "fnb_main_product").First(&mainProduct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("F&B main product not found or not of type fnb_main_product")
		}
		log.Printf("Error finding F&B main product: %v", err)
		return nil, errors.New("failed to retrieve F&B main product")
	}

	if len(mainProduct.Recipes) == 0 {
		return nil, errors.New("F&B main product has no recipes defined")
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deduct component stocks
	for _, recipe := range mainProduct.Recipes {
		if recipe.Component.ID == 0 { // Ensure component is loaded
			tx.Rollback()
			return nil, errors.New("recipe component not found")
		}
		requiredQuantity := recipe.Quantity * req.QuantityToProduce

		// Deduct stock for the component
		err := s.DeductStockForSale(tx, outlet.ID, &recipe.Component.ID, nil, requiredQuantity, userID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error deducting stock for component %s: %v", recipe.Component.Name, err)
			return nil, errors.New("failed to deduct component stock: " + err.Error())
		}
	}

	// Optionally, increase stock of the F&B main product itself
	// This assumes you want to track the stock of the finished F&B product.
	// If not, you can remove this block.
	var mainProductStock models.Stock
	err = tx.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outlet.ID, mainProduct.ID, ownerID).First(&mainProductStock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new stock entry for the main product
			mainProductStock = models.Stock{
				OutletID:  outlet.ID,
				ProductID: &mainProduct.ID,
				Quantity:  req.QuantityToProduce,
				UserID:    ownerID,
			}
			if err := tx.Create(&mainProductStock).Error; err != nil {
				log.Printf("Error creating stock for main F&B product: %v", err)
				return nil, errors.New("failed to create stock for main F&B product")
			}
		} else {
			log.Printf("Error finding stock for main F&B product: %v", err)
			return nil, errors.New("failed to retrieve stock for main F&B product")
		}
	} else {
		// Update existing stock for the main product
		mainProductStock.Quantity += req.QuantityToProduce
		if err := tx.Save(&mainProductStock).Error; err != nil {
			log.Printf("Error updating stock for main F&B product: %v", err)
			return nil, errors.New("failed to update stock for main F&B product")
		}
	}

	// Record stock movement for the main product production
	movement := &models.StockMovement{
		OutletID:       outlet.ID,
		ProductID:      &mainProduct.ID,
		QuantityChange: int(req.QuantityToProduce),
		MovementType:   "Production",
		Description:    stringPtr("Produced F&B main product"),
	}
	if err := s.StockMovementService.CreateStockMovementWithTx(tx, movement); err != nil {
		log.Printf("Error recording production stock movement: %v", err)
		// Don't rollback for movement logging failure, but log it.
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing F&B production transaction: %v", err)
		return nil, errors.New("failed to complete F&B production")
	}

	return &dtos.FNBProductionResponse{
		FNBMainProductUuid: mainProduct.Uuid,
		ProductName:        mainProduct.Name,
		QuantityProduced:   req.QuantityToProduce,
		Message:            "F&B main product produced successfully, components deducted.",
	}, nil
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
func (s *StockService) AddStockFromSale(tx *gorm.DB, outletID uint, productID *uint, productVariantID *uint, quantity float64, userID uint) error {
	var stock models.Stock
	var query *gorm.DB

	if productVariantID != nil {
		query = tx.Where("outlet_id = ? AND product_variant_id = ? AND user_id = ?", outletID, *productVariantID, userID)
	} else if productID != nil {
		query = tx.Where("outlet_id = ? AND product_id = ? AND user_id = ?", outletID, *productID, userID)
	} else {
		return errors.New("product_id or product_variant_id is required for stock addition")
	}

	if err := query.First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If stock not found, create a new entry
			stock = models.Stock{
				OutletID: outletID,
				Quantity: quantity,
				UserID:   userID,
			}
			if productVariantID != nil {
				stock.ProductVariantID = productVariantID
			} else {
				stock.ProductID = productID
			}
			if err := tx.Create(&stock).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		stock.Quantity += quantity
		if err := tx.Save(&stock).Error; err != nil {
			return err
		}
	}

	// Record stock movement
	movement := &models.StockMovement{
		OutletID:         outletID,
		ProductID:        productID,
		ProductVariantID: productVariantID,
		QuantityChange:   int(quantity),
		MovementType:     "Return",
		Description:      stringPtr("Addition from sale return/correction"),
	}
	return s.StockMovementService.CreateStockMovementWithTx(tx, movement)
}

// stringPtr is a helper function to return a pointer to a string.
func stringPtr(s string) *string {
	return &s
}