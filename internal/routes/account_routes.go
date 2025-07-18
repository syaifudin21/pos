package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/handlers"
	internalmw "github.com/msyaifudin/pos/internal/middleware" // Keep this, as SelfAuthorize is used
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/validators"
	"gorm.io/gorm" // Import gorm
)

func RegisterAccountRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize services and handlers for account routes
	userContextService := services.NewUserContextService(db)
	authService := services.NewAuthService(db) // Assuming authService is needed for authHandler
	authHandler := handlers.NewAuthHandler(authService, userContextService)

	ipaymuService := services.NewIpaymuService(db, userContextService)
	ipaymuHandler := handlers.NewIpaymuHandler(ipaymuService, userContextService)

	userPaymentService := services.NewUserPaymentService(db, userContextService) // Assuming this is needed for tsmService
	tsmLogService := services.NewTsmLogService(db)
	orderPaymentService := services.NewOrderPaymentService(db, userContextService, ipaymuService, nil)
	tsmService := services.NewTsmService(db, userContextService, userPaymentService, tsmLogService, orderPaymentService)
	tsmHandler := handlers.NewTsmHandler(tsmService, userContextService, userPaymentService)

	selfAuthGroup := e.Group("")
	selfAuthGroup.Use(internalmw.SelfAuthorize())
	{
		accountGroup := selfAuthGroup.Group("/account")
		accountGroup.GET("/profile", authHandler.GetProfile)
		accountGroup.PUT("/password", authHandler.UpdatePassword, WithValidation(&dtos.UpdatePasswordRequest{}, validators.ValidateUpdatePasswordRequest))
		accountGroup.POST("/email/otp", authHandler.SendOTPForEmailUpdate, WithValidation(&dtos.SendOTPRequest{}, validators.ValidateSendOTPRequest))
		accountGroup.PUT("/email", authHandler.UpdateEmail, WithValidation(&dtos.UpdateEmailRequest{}, validators.ValidateUpdateEmailRequest))

		paymentGroup := selfAuthGroup.Group("/ipaymu")
		paymentGroup.POST("/register", ipaymuHandler.RegisterIpaymu, WithValidation(&dtos.RegisterIpaymuRequest{}, validators.ValidateRegisterIpaymu))
		paymentGroup.POST("/direct-payment", ipaymuHandler.CreateDirectPayment, WithValidation(&dtos.CreateDirectPaymentRequest{}, validators.ValidateCreateDirectPayment))

		tsmGroup := selfAuthGroup.Group("/tsm")
		tsmGroup.POST("/register", tsmHandler.RegisterTsm, WithValidation(&dtos.TsmRegisterRequest{}, validators.ValidateTsmRegister))
	}
}
