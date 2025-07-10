package services

import (
	"context"
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
}

func NewOrderService(db *gorm.DB, stockService *StockService, ipaymuService *IpaymuService, userContextService *UserContextService) *OrderService {
	return &OrderService{DB: db, StockService: stockService, IpaymuService: ipaymuService, UserContextService: userContextService}
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
		OutletID:      outlet.ID,
		UserID:        ownerID,
		Status:        "completed",
		TotalAmount:   0,
		PaymentMethod: req.PaymentMethod,
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

		if item.ProductVariantUuid != uuid.Nil {
			if err := tx.Where("uuid = ? AND user_id = ?", item.ProductVariantUuid, ownerID).First(&variant).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("product variant not found")
			}
			price = variant.Price
			variantID = &variant.ID
		} else if item.ProductUuid != uuid.Nil {
			if err := tx.Where("uuid = ? AND user_id = ?", item.ProductUuid, ownerID).First(&product).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("product not found")
			}
			price = product.Price
			productID = &product.ID
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

	return &dtos.OrderResponse{
		ID:            order.ID,
		Uuid:          order.Uuid,
		OutletID:      order.OutletID,
		OutletUuid:    outlet.Uuid,
		UserID:        order.UserID,
		UserUuid:      user.Uuid,
		OrderDate:     order.CreatedAt.Format(time.RFC3339),
		TotalAmount:   order.TotalAmount,
		PaymentMethod: order.PaymentMethod,
		Status:        order.Status,
	}, nil
}

// GetOrder retrieves an order by its Uuid.
func (s *OrderService) GetOrderByUuid(uuid uuid.UUID, userID uint) (*dtos.OrderResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var order models.Order
	if err := s.DB.Preload("Outlet").Preload("User").Where("uuid = ? AND user_id = ?", uuid, ownerID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		log.Printf("Error getting order by uuid: %v", err)
		return nil, errors.New("failed to retrieve order")
	}

	return &dtos.OrderResponse{
		ID:            order.ID,
		Uuid:          order.Uuid,
		OutletID:      order.OutletID,
		OutletUuid:    order.Outlet.Uuid,
		UserID:        order.UserID,
		UserUuid:      order.User.Uuid,
		OrderDate:     order.CreatedAt.Format(time.RFC3339),
		TotalAmount:   order.TotalAmount,
		PaymentMethod: order.PaymentMethod,
		Status:        order.Status,
	}, nil
}

// GetOrdersByOutlet retrieves all orders for a specific outlet.
func (s *OrderService) GetOrdersByOutlet(outletUuid uuid.UUID, userID uint) ([]dtos.OrderResponse, error) {
	ownerID, err := s.UserContextService.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, ownerID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	if err := s.DB.Preload("User").Where("outlet_id = ? AND user_id = ?", outlet.ID, ownerID).Find(&orders).Error; err != nil {
		log.Printf("Error getting orders by outlet: %v", err)
		return nil, errors.New("failed to retrieve orders")
	}

	var orderResponses []dtos.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, dtos.OrderResponse{
			ID:            order.ID,
			Uuid:          order.Uuid,
			OutletID:      order.OutletID,
			OutletUuid:    outlet.Uuid, // Use the fetched outlet's UUID
			UserID:        order.UserID,
			UserUuid:      order.User.Uuid,
			OrderDate:     order.CreatedAt.Format(time.RFC3339),
			TotalAmount:   order.TotalAmount,
			PaymentMethod: order.PaymentMethod,
			Status:        order.Status,
		})
	}
	return orderResponses, nil
}
