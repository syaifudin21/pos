package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/handlers"
	"github.com/msyaifudin/pos/internal/services"
	"gorm.io/gorm" // Import gorm
)

func RegisterPublicRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize services and handlers for public routes
	userContextService := services.NewUserContextService(db) // Assuming this is needed for ipaymuService
	ipaymuService := services.NewIpaymuService(db, userContextService)
	ipaymuHandler := handlers.NewIpaymuHandler(ipaymuService, userContextService)

	e.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to POS API!")
	})
	e.POST("/api/payment/ipaymu/notify", ipaymuHandler.IpaymuNotify)
}
