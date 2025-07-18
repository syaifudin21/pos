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
	orderPaymentService := services.NewOrderPaymentService(db, userContextService, nil, nil)
	tsmService := services.NewTsmService(db, userContextService, userPaymentService, tsmLogService, orderPaymentService)

	// Initialize IpaymuService
	ipaymuService := services.NewIpaymuService(db, userContextService)

	// Set IpaymuService in OrderPaymentService
	orderPaymentService.IpaymuService = ipaymuService
	orderPaymentService.TsmService = tsmService

	// Set OrderPaymentService in IpaymuService to resolve circular dependency
	ipaymuService.SetOrderPaymentService(orderPaymentService)

	ipaymuHandler := handlers.NewIpaymuHandler(ipaymuService, userContextService)
	tsmHandler := handlers.NewTsmHandler(tsmService, userContextService, userPaymentService)

	e.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to POS API!")
	})
	e.POST("/api/payment/ipaymu/notify", ipaymuHandler.IpaymuNotify, WithValidation(&dtos.IpaymuNotifyRequest{}, validators.ValidateIpaymuNotify))
	e.POST("/api/payment/tsm/callback", tsmHandler.Callback)
}
