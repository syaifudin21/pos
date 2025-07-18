package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type OrderService struct {
	DB                 *gorm.DB
	StockService       *StockService
	IpaymuService      *IpaymuService
	UserContextService *UserContextService
	UserPaymentService *UserPaymentService
}

func NewOrderService(db *gorm.DB, stockService *StockService, ipaymuService *IpaymuService, userContextService *UserContextService) *OrderService {
	return &OrderService{DB: db, StockService: stockService, IpaymuService: ipaymuService, UserContextService: userContextService, UserPaymentService: NewUserPaymentService(db, userContextService)}
}

func (s *OrderService) CreateOrder(req dtos.CreateOrderRequest, userID uint) (*dtos.OrderResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", req.OutletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order := models.Order{
		OutletID:    outlet.ID,
		UserID:      ownerID,
		Status:      "pending",
		TotalAmount: 0,
	}

	if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create order")
	}

	totalAmount := 0.0
	for _, item := range req.Items {
		var product *models.Product
		var variant *models.ProductVariant
		var price float64
		var productID *uint
		var variantID *uint
		var productName string

		if item.ProductVariantUuid != uuid.Nil {
			if err := tx.Where("uuid = ? AND user_id = ?", item.ProductVariantUuid, ownerID).First(&variant).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("product variant not found")
			}
			price = variant.Price
			variantID = &variant.ID
			productName = variant.Name // Use variant name
		} else if item.ProductUuid != uuid.Nil {
			if err := tx.Where("uuid = ? AND user_id = ?", item.ProductUuid, ownerID).First(&product).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("product not found")
			}
			price = product.Price
			productID = &product.ID
			productName = product.Name // Use product name
		} else {
			tx.Rollback()
			return nil, errors.New("product_uuid or product_variant_uuid is required for each item")
		}

		if err := s.StockService.DeductStockForSale(tx, outlet.ID, productID, variantID, float64(item.Quantity), ownerID); err != nil {
			tx.Rollback()
			return nil, err
		}

		orderItem := models.OrderItem{
			OrderID:          order.ID,
			ProductID:        productID,
			ProductVariantID: variantID,
			Quantity:         float64(item.Quantity),
			Price:            price,
			ProductName:      productName,
		}

		if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to create order item")
		}

		// Process add-ons for the current order item
		for _, addOnReq := range item.AddOns {
			var addOnProduct models.Product
			if err := tx.Where("uuid = ? AND user_id = ? AND type = ?", addOnReq.AddOnUuid, ownerID, "add_on").First(&addOnProduct).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("add-on product not found or not of type add_on")
			}

			orderItemAddOn := models.OrderItemAddOn{
				OrderItemID: orderItem.ID,
				AddOnID:     addOnProduct.ID,
				Quantity:    float64(addOnReq.Quantity),
				Price:       addOnProduct.Price, // Use the add-on's base price
				UserID:      ownerID,
			}

			if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&orderItemAddOn).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("failed to create order item add-on")
			}
			totalAmount += addOnProduct.Price * float64(addOnReq.Quantity)
		}

		totalAmount += price * float64(item.Quantity)
	}

	order.TotalAmount = totalAmount
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update order total")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to commit order transaction")
	}

	// Reload the order with all its relations for the comprehensive response using the main DB connection
	if err := s.DB.Preload("User").Preload("Outlet").Preload("OrderPayments.PaymentMethod").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").Preload("OrderItems.AddOns.AddOn").First(&order, order.ID).Error; err != nil {
		log.Printf("Error preloading order relations after commit: %v", err)
		return nil, errors.New("failed to retrieve full order details after commit")
	}

	return mapOrderToOrderResponse(order, outlet), nil
}

// GetOrder retrieves an order by its Uuid.
func (s *OrderService) GetOrderByUuid(uuid uuid.UUID, userID uint) (*dtos.OrderResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var order models.Order
	if err := s.DB.Preload("User").Preload("Outlet").Preload("OrderPayments.PaymentMethod").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant.Product").Preload("OrderItems.AddOns.AddOn").Preload("OrderItems.OrderPaymentItems.OrderPayment").Where("uuid = ? AND user_id = ?", uuid, ownerID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		log.Printf("Error getting order by uuid: %v", err)
		return nil, errors.New("failed to retrieve order")
	}

	return mapOrderToOrderResponse(order, order.Outlet), nil
}

// GetOrdersByOutlet retrieves all orders for a specific outlet.
func (s *OrderService) GetOrdersByOutlet(outletUuid uuid.UUID, userID uint, status string) ([]dtos.SimpleOrderResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	query := s.DB.Preload("User").Where("outlet_id = ? AND user_id = ?", outlet.ID, ownerID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&orders).Error; err != nil {
		log.Printf("Error getting orders by outlet: %v", err)
		return nil, errors.New("failed to retrieve orders")
	}

	var orderResponses []dtos.SimpleOrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, *mapOrderToSimpleOrderResponse(order))
	}
	return orderResponses, nil
}

func (s *OrderService) UpdateOrderItem(orderUuid uuid.UUID, req dtos.UpdateOrderItemRequest, userID uint) (*dtos.OrderResponse, error) {
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

	var order models.Order
	if err := tx.Where("uuid = ? AND user_id = ?", orderUuid, ownerID).First(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("order not found")
	}

	var orderItem models.OrderItem
	if err := tx.Preload("AddOns").Preload("OrderPaymentItems.OrderPayment").Where("uuid = ? AND order_id = ?", req.OrderItemUuid, order.ID).First(&orderItem).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("order item not found")
	}

	// Check if the order item is already paid
	for _, opItem := range orderItem.OrderPaymentItems {
		if opItem.OrderPayment != nil && opItem.OrderPayment.IsPaid {
			tx.Rollback()
			return nil, errors.New("cannot update a paid order item")
		}
	}

	// Return stock for old item and add-ons
	if orderItem.ProductID != nil {
		if err := s.StockService.AddStockFromSale(tx, order.OutletID, orderItem.ProductID, nil, orderItem.Quantity, ownerID); err != nil {
			tx.Rollback()
			return nil, err
		}
	} else if orderItem.ProductVariantID != nil {
		if err := s.StockService.AddStockFromSale(tx, order.OutletID, nil, orderItem.ProductVariantID, orderItem.Quantity, ownerID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, addOn := range orderItem.AddOns {
		if err := s.StockService.AddStockFromSale(tx, order.OutletID, &addOn.AddOnID, nil, addOn.Quantity, ownerID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Delete old add-ons
	if err := tx.Where("order_item_id = ?", orderItem.ID).Delete(&models.OrderItemAddOn{}).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to delete old order item add-ons")
	}

	var product *models.Product
	var variant *models.ProductVariant
	var price float64
	var productID *uint
	var variantID *uint
	var productName string

	if req.ProductVariantUuid != uuid.Nil {
		if err := tx.Where("uuid = ? AND user_id = ?", req.ProductVariantUuid, ownerID).First(&variant).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product variant not found")
		}
		price = variant.Price
		variantID = &variant.ID
		productName = variant.Name
	} else if req.ProductUuid != uuid.Nil {
		if err := tx.Where("uuid = ? AND user_id = ?", req.ProductUuid, ownerID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}
		price = product.Price
		productID = &product.ID
		productName = product.Name
	} else {
		tx.Rollback()
		return nil, errors.New("product_uuid or product_variant_uuid is required for each item")
	}

	if err := s.StockService.DeductStockForSale(tx, order.OutletID, productID, variantID, float64(req.Quantity), ownerID); err != nil {
		tx.Rollback()
		return nil, err
	}

	orderItem.ProductID = productID
	orderItem.ProductVariantID = variantID
	orderItem.Quantity = float64(req.Quantity)
	orderItem.Price = price
	orderItem.ProductName = productName

	if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&orderItem).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update order item")
	}

	// Process new add-ons
	for _, addOnReq := range req.AddOns {
		var addOnProduct models.Product
		if err := tx.Where("uuid = ? AND user_id = ? AND type = ?", addOnReq.AddOnUuid, ownerID, "add_on").First(&addOnProduct).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("add-on product not found or not of type add_on")
		}

		orderItemAddOn := models.OrderItemAddOn{
			OrderItemID: orderItem.ID,
			AddOnID:     addOnProduct.ID,
			Quantity:    float64(addOnReq.Quantity),
			Price:       addOnProduct.Price,
			UserID:      ownerID,
		}

		if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&orderItemAddOn).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to create order item add-on")
		}
	}

	// Recalculate total amount for the order
	if err := s.recalculateOrderTotal(tx, &order, ownerID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to commit order item update transaction")
	}

	// Reload the order with all its relations for the comprehensive response
	if err := s.DB.Preload("User").Preload("Outlet").Preload("OrderPayments.PaymentMethod").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").Preload("OrderItems.AddOns.AddOn").Preload("OrderItems.OrderPaymentItems.OrderPayment").First(&order, order.ID).Error; err != nil {
		log.Printf("Error preloading order relations after commit: %v", err)
		return nil, errors.New("failed to retrieve full order details after commit")
	}

	return mapOrderToOrderResponse(order, order.Outlet), nil
}

func (s *OrderService) DeleteOrderItem(orderUuid uuid.UUID, req dtos.DeleteOrderItemRequest, userID uint) (*dtos.OrderResponse, error) {
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

	var order models.Order
	if err := tx.Where("uuid = ? AND user_id = ?", orderUuid, ownerID).First(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("order not found")
	}

	var orderItem models.OrderItem
	if err := tx.Preload("AddOns").Preload("OrderPaymentItems.OrderPayment").Where("uuid = ? AND order_id = ?", req.OrderItemUuid, order.ID).First(&orderItem).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("order item not found")
	}
	log.Printf("Retrieved OrderItem for deletion: ID=%d, UUID=%s", orderItem.ID, orderItem.Uuid.String())

	if orderItem.ID == 0 {
		tx.Rollback()
		return nil, errors.New("retrieved order item has invalid ID for deletion")
	}

	// Check if the order item is already paid
	for _, opItem := range orderItem.OrderPaymentItems {
		if opItem.OrderPayment != nil && opItem.OrderPayment.IsPaid {
			tx.Rollback()
			return nil, errors.New("cannot delete a paid order item")
		}
	}

	// Return stock for the deleted item and its add-ons
	if orderItem.ProductID != nil {
		if err := s.StockService.AddStockFromSale(tx, order.OutletID, orderItem.ProductID, nil, orderItem.Quantity, ownerID); err != nil {
			tx.Rollback()
			return nil, err
		}
	} else if orderItem.ProductVariantID != nil {
		if err := s.StockService.AddStockFromSale(tx, order.OutletID, nil, orderItem.ProductVariantID, orderItem.Quantity, ownerID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, addOn := range orderItem.AddOns {
		if err := s.StockService.AddStockFromSale(tx, order.OutletID, &addOn.AddOnID, nil, addOn.Quantity, ownerID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Delete the order item and its add-ons
	log.Printf("Attempting to delete OrderItemAddOn for order_item_id: %d", orderItem.ID)
	result := tx.Unscoped().Where("order_item_id = ?", orderItem.ID).Delete(&models.OrderItemAddOn{})
	if result.Error != nil {
		tx.Rollback()
		return nil, errors.New("failed to delete order item add-ons: " + result.Error.Error())
	}
	log.Printf("Deleted %d OrderItemAddOn rows for order_item_id: %d", result.RowsAffected, orderItem.ID)

	// Delete associated order payment items
	log.Printf("Attempting to delete OrderPaymentItem for order_item_id: %d", orderItem.ID)
	result = tx.Unscoped().Where("order_item_id = ?", orderItem.ID).Delete(&models.OrderPaymentItem{})
	if result.Error != nil {
		tx.Rollback()
		return nil, errors.New("failed to delete associated order payment items: " + result.Error.Error())
	}
	log.Printf("Deleted %d OrderPaymentItem rows for order_item_id: %d", result.RowsAffected, orderItem.ID)

		log.Printf("Attempting to delete OrderItem with ID: %d (UUID: %s)", orderItem.ID, orderItem.Uuid.String())
	result = tx.Unscoped().Delete(&models.OrderItem{}, orderItem.ID) // Reverted to Unscoped().Delete()
	if result.Error != nil {
		tx.Rollback()
		return nil, errors.New("failed to delete order item: " + result.Error.Error())
	}
	log.Printf("Deleted %d OrderItem rows for ID: %d", result.RowsAffected, orderItem.ID)

	// Add verification step here
	var verifyItem models.OrderItem
	checkErr := tx.Unscoped().Where("id = ?", orderItem.ID).First(&verifyItem).Error
	if checkErr == nil {
		tx.Rollback()
		return nil, errors.New("order item still exists after deletion attempt")
	}
	if !errors.Is(checkErr, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, errors.New("error verifying order item deletion: " + checkErr.Error())
	}
	log.Printf("Order item with ID %d successfully verified as deleted from DB.", orderItem.ID)

	// Recalculate total amount for the order
	if err := s.recalculateOrderTotal(tx, &order, ownerID); err != nil {
		tx.Rollback()
		return nil, err
	}

	log.Printf("Attempting to commit transaction for order item deletion.")
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		tx.Rollback()
		return nil, errors.New("failed to commit order item deletion transaction: " + err.Error())
	}
	log.Printf("Transaction committed successfully.")

	// Fetch a fresh order object after commit
	var freshOrder models.Order
	if err := s.DB.Preload("User").Preload("Outlet").Preload("OrderPayments.PaymentMethod").Where("uuid = ? AND user_id = ?", orderUuid, ownerID).First(&freshOrder).Error; err != nil {
		log.Printf("Error fetching fresh order after commit: %v", err)
		return nil, errors.New("failed to retrieve fresh order details after commit")
	}

	// Explicitly re-fetch OrderItems to ensure deleted item is not present
	// Use a new DB session to bypass any potential transaction-level caching
	if err := s.DB.Session(&gorm.Session{NewDB: true}).Where("order_id = ?", freshOrder.ID).Find(&freshOrder.OrderItems).Error; err != nil {
		log.Printf("Error re-fetching order items after deletion: %v", err)
		return nil, errors.New("failed to re-fetch order items after deletion")
	}

	// Now preload the necessary relations for the re-fetched OrderItems
	for i := range freshOrder.OrderItems {
		if freshOrder.OrderItems[i].ProductID != nil {
			s.DB.Where("id = ?", freshOrder.OrderItems[i].ProductID).First(&freshOrder.OrderItems[i].Product)
		}
		if freshOrder.OrderItems[i].ProductVariantID != nil {
			s.DB.Preload("Product").Where("id = ?", freshOrder.OrderItems[i].ProductVariantID).First(&freshOrder.OrderItems[i].ProductVariant)
		}
		s.DB.Preload("AddOn").Where("order_item_id = ?", freshOrder.OrderItems[i].ID).Find(&freshOrder.OrderItems[i].AddOns)
		s.DB.Preload("OrderPayment").Where("order_item_id = ?", freshOrder.OrderItems[i].ID).Find(&freshOrder.OrderItems[i].OrderPaymentItems)
	}

	log.Printf("Fresh Order %d now has %d items after re-fetch.", freshOrder.ID, len(freshOrder.OrderItems))
	log.Printf("Final OrderItems count before mapping to response: %d", len(freshOrder.OrderItems))
	for i, item := range freshOrder.OrderItems {
		log.Printf("Item %d UUID: %s", i, item.Uuid.String())
	}

	return mapOrderToOrderResponse(freshOrder, freshOrder.Outlet), nil
}

func (s *OrderService) CreateOrderItem(orderUuid uuid.UUID, req dtos.CreateOrderItemRequest, userID uint) (*dtos.OrderResponse, error) {
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

	var order models.Order
	if err := tx.Preload("OrderPayments").Where("uuid = ? AND user_id = ?", orderUuid, ownerID).First(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("order not found")
	}

	// Check if the order is already paid
	if order.Status == "paid" {
		tx.Rollback()
		return nil, errors.New("cannot add item to a paid order")
	}

	var product *models.Product
	var variant *models.ProductVariant
	var price float64
	var productID *uint
	var variantID *uint
	var productName string

	if req.ProductVariantUuid != uuid.Nil {
		if err := tx.Where("uuid = ? AND user_id = ?", req.ProductVariantUuid, ownerID).First(&variant).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product variant not found")
		}
		price = variant.Price
		variantID = &variant.ID
		productName = variant.Name
	} else if req.ProductUuid != uuid.Nil {
		if err := tx.Where("uuid = ? AND user_id = ?", req.ProductUuid, ownerID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}
		price = product.Price
		productID = &product.ID
		productName = product.Name
	} else {
		tx.Rollback()
		return nil, errors.New("product_uuid or product_variant_uuid is required for each item")
	}

	if err := s.StockService.DeductStockForSale(tx, order.OutletID, productID, variantID, float64(req.Quantity), ownerID); err != nil {
		tx.Rollback()
		return nil, err
	}

	orderItem := models.OrderItem{
		OrderID:          order.ID,
		ProductID:        productID,
		ProductVariantID: variantID,
		Quantity:         float64(req.Quantity),
		Price:            price,
		ProductName:      productName,
	}

	if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&orderItem).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create order item")
	}

	// Process add-ons for the current order item
	for _, addOnReq := range req.AddOns {
		var addOnProduct models.Product
		if err := tx.Where("uuid = ? AND user_id = ? AND type = ?", addOnReq.AddOnUuid, ownerID, "add_on").First(&addOnProduct).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("add-on product not found or not of type add_on")
		}

		orderItemAddOn := models.OrderItemAddOn{
			OrderItemID: orderItem.ID,
			AddOnID:     addOnProduct.ID,
			Quantity:    float64(addOnReq.Quantity),
			Price:       addOnProduct.Price,
			UserID:      ownerID,
		}

		if err := tx.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&orderItemAddOn).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to create order item add-on")
		}
	}

	// Recalculate total amount for the order
	if err := s.recalculateOrderTotal(tx, &order, ownerID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to commit order item creation transaction")
	}

	// Reload the order with all its relations for the comprehensive response
	if err := s.DB.Preload("User").Preload("Outlet").Preload("OrderPayments.PaymentMethod").Preload("OrderItems.Product").Preload("OrderItems.ProductVariant").Preload("OrderItems.AddOns.AddOn").Preload("OrderItems.OrderPaymentItems.OrderPayment").First(&order, order.ID).Error; err != nil {
		log.Printf("Error preloading order relations after commit: %v", err)
		return nil, errors.New("failed to retrieve full order details after commit")
	}

	return mapOrderToOrderResponse(order, order.Outlet), nil
}

func (s *OrderService) recalculateOrderTotal(tx *gorm.DB, order *models.Order, ownerID uint) error {
	var orderItems []models.OrderItem
	if err := tx.Preload("AddOns").Where("order_id = ?", order.ID).Find(&orderItems).Error; err != nil {
		return errors.New("failed to retrieve order items for recalculation")
	}

	newTotalAmount := 0.0
	for _, item := range orderItems {
		newTotalAmount += item.Price * item.Quantity
		for _, addOn := range item.AddOns {
			newTotalAmount += addOn.Price * addOn.Quantity
		}
	}

	order.TotalAmount = newTotalAmount
	if err := tx.Save(order).Error; err != nil {
		return errors.New("failed to update order total amount")
	}
	return nil
}

func mapOrderToSimpleOrderResponse(order models.Order) *dtos.SimpleOrderResponse {
	return &dtos.SimpleOrderResponse{
		Uuid:        order.Uuid,
		OrderDate:   order.CreatedAt.Format(time.RFC3339),
		TotalAmount: order.TotalAmount,
		PaidAmount:  order.PaidAmount,
		Status:      order.Status,
	}
}

func mapOrderToOrderResponse(order models.Order, outlet models.Outlet) *dtos.OrderResponse {
	var orderItemsResponse []dtos.OrderItemDetailResponse
	for _, item := range order.OrderItems {
		var productName string
		var productUuid = uuid.Nil
		var productVariantUuid = uuid.Nil

		if item.Product != nil {
			productName = item.Product.Name
			productUuid = item.Product.Uuid
		}
		if item.ProductVariant != nil {
			productName = item.ProductVariant.Name
			productVariantUuid = item.ProductVariant.Uuid
			if item.ProductVariant.Product.Uuid != uuid.Nil { // Get parent product UUID if variant exists
				productUuid = item.ProductVariant.Product.Uuid
			}
		}

		var addOnsResponse []dtos.OrderItemAddonDetailResponse
		for _, addOn := range item.AddOns {
			if addOn.AddOn.ID != 0 { // Check if AddOn relation is loaded
				addOnsResponse = append(addOnsResponse, dtos.OrderItemAddonDetailResponse{
					Uuid:     addOn.AddOn.Uuid,
					Name:     addOn.AddOn.Name,
					Quantity: int(addOn.Quantity),
				})
			}
		}

		itemPrice := item.Price
		itemTotal := item.Price * item.Quantity
		for _, addOn := range item.AddOns {
			itemTotal += addOn.Price * addOn.Quantity
		}

		var itemIsPaid bool
		for _, opItem := range item.OrderPaymentItems {
			if opItem.OrderPayment != nil && opItem.OrderPayment.IsPaid {
				itemIsPaid = true
				break
			}
		}

		orderItemsResponse = append(orderItemsResponse, dtos.OrderItemDetailResponse{
			ID:                 item.ID,
			Uuid:               item.Uuid,
			ProductUuid:        productUuid,
			ProductVariantUuid: productVariantUuid,
			Name:               productName,
			Quantity:           int(item.Quantity),
			Price:              itemPrice,
			Total:              itemTotal,
			IsPaid:             itemIsPaid,
			AddOns:             addOnsResponse,
		})
	}

	var paymentsResponse []dtos.OrderPaymentDetailResponse
	var paymentMethods []string
	uniquePaymentMethods := make(map[string]bool)

	for _, payment := range order.OrderPayments {
		if payment.PaymentMethod.ID != 0 {
			paymentsResponse = append(paymentsResponse, dtos.OrderPaymentDetailResponse{
				Uuid:            payment.Uuid,
				PaymentMethodID: payment.PaymentMethodID,
				PaidAmount:      payment.AmountPaid,
				CustomerName:    payment.CustomerName,
				CustomerEmail:   payment.CustomerEmail,
				CustomerPhone:   payment.CustomerPhone,
				Name:            payment.PaymentMethod.Name,
				PaymentMethod:   payment.PaymentMethod.PaymentMethod,
				PaymentChannel:  payment.PaymentMethod.PaymentChannel,
				ChangeAmount:    payment.ChangeAmount,
				IsPaid:          payment.IsPaid,
				ReferenceID:     payment.ReferenceID,
				CreatedAt:       payment.CreatedAt.Format(time.RFC3339),
				PaidAt:          payment.PaidAt,
				Extra: func() interface{} {
					var extraData interface{}
					if payment.Extra != "" {
						err := json.Unmarshal([]byte(payment.Extra), &extraData)
						if err != nil {
							log.Printf("Failed to unmarshal Extra field for OrderPaymentDetailResponse: %v", err)
							return nil // Return nil if unmarshaling fails
						}
					}
					return extraData
				}(),
			})

			if payment.PaymentMethod.PaymentChannel != "" {
				uniquePaymentMethods[payment.PaymentMethod.PaymentChannel] = true
			}
		}
	}

	for method := range uniquePaymentMethods {
		paymentMethods = append(paymentMethods, method)
	}

	var createdBy *dtos.UserDetailResponse
	if order.User.ID != 0 {
		createdBy = &dtos.UserDetailResponse{
			Uuid: order.User.Uuid,
			Name: order.User.Name,
		}
	}

	return &dtos.OrderResponse{
		Uuid:           order.Uuid,
		OrderDate:      order.CreatedAt.Format(time.RFC3339),
		TotalAmount:    order.TotalAmount,
		PaidAmount:     order.PaidAmount,
		Status:         order.Status,
		PaymentMethods: paymentMethods,
		CreatedBy:      createdBy,
		Outlet: dtos.OutletDetailResponse{
			Uuid:    outlet.Uuid,
			Name:    outlet.Name,
			Address: outlet.Address,
			Contact: outlet.Contact,
		},
		Payments: paymentsResponse,
		Items:    orderItemsResponse,
	}
}
