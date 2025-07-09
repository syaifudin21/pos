package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/handlers"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/validators"
	"gorm.io/gorm" // Import gorm
)

func RegisterAuthRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize services and handlers for auth routes
	userContextService := services.NewUserContextService(db)
	authService := services.NewAuthService(db)
	authHandler := handlers.NewAuthHandler(authService, userContextService)
	googleOAuthService := services.NewGoogleOAuthService(db, authService)
	googleOAuthHandler := handlers.NewGoogleOAuthHandler(googleOAuthService)

	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.RegisterOwner, WithValidation(&dtos.RegisterOwnerRequest{}, validators.ValidateRegisterOwnerRequest))
	authGroup.POST("/verify-otp", authHandler.VerifyOTP, WithValidation(&dtos.VerifyOTPRequest{}, validators.ValidateVerifyOTPRequest))
	authGroup.POST("/login", authHandler.Login, WithValidation(&dtos.LoginRequest{}, validators.ValidateLoginRequest))
	authGroup.POST("/forgot-password", authHandler.ForgotPassword, WithValidation(&dtos.ForgotPasswordRequest{}, validators.ValidateForgotPasswordRequest))
	authGroup.POST("/reset-password", authHandler.ResetPassword, WithValidation(&dtos.ResetPasswordRequest{}, validators.ValidateResetPasswordRequest))
	authGroup.POST("/resend-verification-email", authHandler.ResendVerificationEmail, WithValidation(&dtos.ResendEmailRequest{}, validators.ValidateResendEmailRequest))
	authGroup.GET("/google/login", googleOAuthHandler.GoogleLogin)
	authGroup.GET("/google/callback", googleOAuthHandler.GoogleCallback)
}
