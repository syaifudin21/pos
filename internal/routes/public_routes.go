package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/handlers"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/validators"
	"gorm.io/gorm" // Import gorm
)

func RegisterPublicRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize services and handlers for public routes
	userContextService := services.NewUserContextService(db)

	// Dependencies for TsmService & OrderPaymentService
	userPaymentService := services.NewUserPaymentService(db, userContextService)
	tsmLogService := services.NewTsmLogService(db)
	tsmService := services.NewTsmService(db, userContextService, userPaymentService, tsmLogService)

	// Initialize OrderPaymentService first
	orderPaymentService := services.NewOrderPaymentService(db, userContextService, nil, tsmService) // ipaymuService is nil for now to avoid circular dependency

	// Initialize IpaymuService
	ipaymuService := services.NewIpaymuService(db, userContextService)

	// Set IpaymuService in OrderPaymentService
	orderPaymentService.IpaymuService = ipaymuService

	// Set OrderPaymentService in IpaymuService to resolve circular dependency
	ipaymuService.SetOrderPaymentService(orderPaymentService)

	ipaymuHandler := handlers.NewIpaymuHandler(ipaymuService, userContextService)

	e.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to POS API!")
	})
	e.POST("/api/payment/ipaymu/notify", ipaymuHandler.IpaymuNotify, WithValidation(&dtos.IpaymuNotifyRequest{}, validators.ValidateIpaymuNotify))
}
