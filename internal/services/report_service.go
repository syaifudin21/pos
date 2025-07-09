package services

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type ReportService struct {
	DB *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{DB: db}
}

// SalesByOutletReport generates a sales report for a specific outlet within a date range.
func (s *ReportService) SalesByOutletReport(outletUuid uuid.UUID, startDate, endDate time.Time, userID uint) ([]models.Order, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, userID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var orders []models.Order
	err := s.DB.Preload("OrderItems.Product").
		Where("outlet_id = ? AND user_id = ? AND created_at BETWEEN ? AND ?", outlet.ID, userID, startDate, endDate.Add(24*time.Hour)).
		Find(&orders).Error

	if err != nil {
		log.Printf("Error generating sales by outlet report: %v", err)
		return nil, errors.New("failed to generate report")
	}

	return orders, nil
}

// SalesByProductReport generates a sales report for a specific product within a date range.
func (s *ReportService) SalesByProductReport(productUuid uuid.UUID, startDate, endDate time.Time, userID uint) ([]models.OrderItem, error) {
	var product models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", productUuid, userID).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}

	var orderItems []models.OrderItem
	err := s.DB.Preload("Order.Outlet").Preload("Order.User").
		Where("product_id = ? AND user_id = ? AND created_at BETWEEN ? AND ?", product.ID, userID, startDate, endDate.Add(24*time.Hour)).
		Find(&orderItems).Error

	if err != nil {
		log.Printf("Error generating sales by product report: %v", err)
		return nil, errors.New("failed to generate report")
	}

	return orderItems, nil
}

// StockReport generates a stock report for a specific outlet.
func (s *ReportService) StockReport(outletUuid uuid.UUID, userID uint) ([]dtos.StockReportResponse, error) {
	var outlet models.Outlet
	if err := s.DB.Where("uuid = ? AND user_id = ?", outletUuid, userID).First(&outlet).Error; err != nil {
		return nil, errors.New("outlet not found")
	}

	var stocks []models.Stock
	err := s.DB.
		Preload("Product").
		Preload("ProductVariant.Product"). // Preload the parent product of the variant
		Where("outlet_id = ? AND user_id = ?", outlet.ID, userID).
		Find(&stocks).Error

	if err != nil {
		log.Printf("Error generating stock report: %v", err)
		return nil, errors.New("failed to generate stock report")
	}

	var report []dtos.StockReportResponse
	for _, stock := range stocks {
		if stock.ProductID != nil && stock.Product != nil {
			report = append(report, dtos.StockReportResponse{
				ProductName: stock.Product.Name,
				ProductSku:  stock.Product.SKU,
				Quantity:    stock.Quantity,
			})
		} else if stock.ProductVariantID != nil && stock.ProductVariant != nil && stock.ProductVariant.Product.ID != 0 {
			report = append(report, dtos.StockReportResponse{
				ProductName: stock.ProductVariant.Product.Name, // Get product name from the preloaded parent
				VariantName: stock.ProductVariant.Name,
				VariantSku:  stock.ProductVariant.SKU,
				Quantity:    stock.Quantity,
			})
		}
	}

	return report, nil
}
