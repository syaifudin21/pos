package services

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
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
