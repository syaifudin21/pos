package routes

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/handlers"
	internalmw "github.com/msyaifudin/pos/internal/middleware"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/validators"
)

func withValidation(dtoType interface{}, validatorFunc interface{}) echo.MiddlewareFunc {
	return internalmw.ValidationMiddleware(dtoType, func(data interface{}) []string {
		// Use reflection to call the actual validator function with the correct type
		validatorValue := reflect.ValueOf(validatorFunc)
		// Ensure data is a pointer if the validator expects a pointer
		var arg reflect.Value
		if reflect.TypeOf(dtoType).Kind() == reflect.Ptr {
			arg = reflect.ValueOf(data)
		} else {
			arg = reflect.ValueOf(data).Elem()
		}

		results := validatorValue.Call([]reflect.Value{arg})
		if len(results) > 0 && !results[0].IsNil() {
			return results[0].Interface().([]string)
		}
		return nil
	})
}

func RegisterRoutes(e *echo.Echo) {
	// Initialize services and handlers
	userContextService := services.NewUserContextService(database.DB)

	// Initialize iPaymu service
	ipaymuService := services.NewIpaymuService(database.DB, userContextService)
	authService := services.NewAuthService(database.DB)
	authHandler := handlers.NewAuthHandler(authService, userContextService)
	googleOAuthService := services.NewGoogleOAuthService(database.DB, authService)
	googleOAuthHandler := handlers.NewGoogleOAuthHandler(googleOAuthService)
	productService := services.NewProductService(database.DB, userContextService)
	productHandler := handlers.NewProductHandler(productService, userContextService)
	outletService := services.NewOutletService(database.DB, userContextService)
	outletHandler := handlers.NewOutletHandler(outletService, userContextService)
	recipeService := services.NewRecipeService(database.DB, userContextService)
	recipeHandler := handlers.NewRecipeHandler(recipeService, userContextService)
	stockService := services.NewStockService(database.DB, userContextService)
	stockHandler := handlers.NewStockHandler(stockService, userContextService)
	orderService := services.NewOrderService(database.DB, stockService, ipaymuService, userContextService)
	orderHandler := handlers.NewOrderHandler(orderService, userContextService)
	reportService := services.NewReportService(database.DB)
	reportHandler := handlers.NewReportHandler(reportService, userContextService)
	supplierService := services.NewSupplierService(database.DB, userContextService)
	supplierHandler := handlers.NewSupplierHandler(supplierService, userContextService)
	poService := services.NewPurchaseOrderService(database.DB, stockService, userContextService)
	poHandler := handlers.NewPurchaseOrderHandler(poService, userContextService)
	ipaymuHandler := handlers.NewIpaymuHandler(ipaymuService, userContextService)
	userPaymentService := services.NewUserPaymentService(database.DB, userContextService)
	userPaymentHandler := handlers.NewUserPaymentHandler(userPaymentService, userContextService)
	tsmService := services.NewTsmService(database.DB, userContextService, userPaymentService)
	tsmHandler := handlers.NewTsmHandler(tsmService, userContextService, userPaymentService)

	// Public routes (no specific middleware)
	e.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to POS API!")
	})
	e.POST("/api/payment/ipaymu/notify", ipaymuHandler.IpaymuNotify)

	// Auth routes (public, but grouped)
	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.RegisterOwner, withValidation(&dtos.RegisterAdminRequest{}, validators.ValidateRegisterAdminRequest))
	authGroup.POST("/verify-otp", authHandler.VerifyOTP, withValidation(&dtos.VerifyOTPRequest{}, validators.ValidateVerifyOTPRequest))
	authGroup.POST("/login", authHandler.Login, withValidation(&dtos.LoginRequest{}, validators.ValidateLoginRequest))
	authGroup.POST("/forgot-password", authHandler.ForgotPassword, withValidation(&dtos.ForgotPasswordRequest{}, validators.ValidateForgotPasswordRequest))
	authGroup.POST("/reset-password", authHandler.ResetPassword, withValidation(&dtos.ResetPasswordRequest{}, validators.ValidateResetPasswordRequest))
	authGroup.POST("/resend-verification-email", authHandler.ResendVerificationEmail, withValidation(&dtos.ResendEmailRequest{}, validators.ValidateResendEmailRequest))
	authGroup.GET("/google/login", googleOAuthHandler.GoogleLogin)
	authGroup.GET("/google/callback", googleOAuthHandler.GoogleCallback)

	// Routes requiring internalmw.SelfAuthorize()
	selfAuthGroup := e.Group("")
	selfAuthGroup.Use(internalmw.SelfAuthorize())
	{
		accountGroup := selfAuthGroup.Group("/account")
		accountGroup.GET("/profile", authHandler.GetProfile)
		accountGroup.PUT("/password", authHandler.UpdatePassword, withValidation(&dtos.UpdatePasswordRequest{}, validators.ValidateUpdatePasswordRequest))
		accountGroup.POST("/email/otp", authHandler.SendOTPForEmailUpdate, withValidation(&dtos.SendOTPRequest{}, validators.ValidateSendOTPRequest))
		accountGroup.PUT("/email", authHandler.UpdateEmail, withValidation(&dtos.UpdateEmailRequest{}, validators.ValidateUpdateEmailRequest))

		paymentGroup := selfAuthGroup.Group("/ipaymu")
		paymentGroup.POST("/register", ipaymuHandler.RegisterIpaymu, withValidation(&dtos.RegisterIpaymuRequest{}, validators.ValidateRegisterIpaymu))
		paymentGroup.POST("/direct-payment", ipaymuHandler.CreateDirectPayment, withValidation(&dtos.CreateDirectPaymentRequest{}, validators.ValidateCreateDirectPayment))

		tsmGroup := selfAuthGroup.Group("/tsm")
		tsmGroup.POST("/register", tsmHandler.RegisterTsm, withValidation(&dtos.TsmRegisterRequest{}, validators.ValidateTsmRegister))
	}

	// Routes requiring internalmw.Authorize()
	authorizedGroup := e.Group("")
	{
		// User management routes (owner only)
		userAdminGroup := authorizedGroup.Group("/users", internalmw.Authorize("users", "manage"))
		userAdminGroup.GET("", authHandler.GetAllUsers)
		userAdminGroup.POST("", authHandler.Register, withValidation(&dtos.RegisterRequest{}, validators.ValidateRegisterRequest))
		userAdminGroup.PUT("/:uuid", authHandler.UpdateUser, withValidation(&dtos.UpdateUserRequest{}, validators.ValidateUpdateUserRequest))
		userAdminGroup.PUT("/:uuid/block", authHandler.BlockUser)
		userAdminGroup.PUT("/:uuid/unblock", authHandler.UnblockUser)
		userAdminGroup.DELETE("/:uuid", authHandler.DeleteUser)

		// User Payment routes (owner only)
		userPaymentGroup := authorizedGroup.Group("/account/payment-methods")
		userPaymentGroup.POST("/activate", userPaymentHandler.ActivateUserPayment, internalmw.Authorize("user_payments", "activate"))
		userPaymentGroup.POST("/deactivate", userPaymentHandler.DeactivateUserPayment, internalmw.Authorize("user_payments", "deactivate"))
		userPaymentGroup.GET("", userPaymentHandler.ListUserPayments, internalmw.Authorize("user_payments", "read"))

		// Product routes
		productGroup := authorizedGroup.Group("/products", internalmw.Authorize("products", "read"))
		productGroup.GET("", productHandler.GetAllProducts)
		productGroup.GET("/:uuid", productHandler.GetProductByID)
		productGroup.POST("", productHandler.CreateProduct, internalmw.Authorize("products", "write"), withValidation(&dtos.ProductCreateRequest{}, validators.ValidateCreateProduct))
		productGroup.PUT("/:uuid", productHandler.UpdateProduct, internalmw.Authorize("products", "write"), withValidation(&dtos.ProductUpdateRequest{}, validators.ValidateUpdateProduct))
		productGroup.DELETE("/:uuid", productHandler.DeleteProduct, internalmw.Authorize("products", "write"))

		outletProductGroup := authorizedGroup.Group("/outlets/:outlet_uuid/products", internalmw.Authorize("products", "read"))
		outletProductGroup.GET("", productHandler.GetProductsByOutlet)

		// Recipe routes
		recipeGroup := authorizedGroup.Group("/recipes", internalmw.Authorize("recipes", "read"))
		recipeGroup.GET("/:uuid", recipeHandler.GetRecipeByUuid)
		recipeGroup.POST("", recipeHandler.CreateRecipe, internalmw.Authorize("recipes", "write"), withValidation(&dtos.CreateRecipeRequest{}, validators.ValidateCreateRecipe))
		recipeGroup.PUT("/:uuid", recipeHandler.UpdateRecipe, internalmw.Authorize("recipes", "write"), withValidation(&dtos.UpdateRecipeRequest{}, validators.ValidateUpdateRecipe))
		recipeGroup.DELETE("/:uuid", recipeHandler.DeleteRecipe, internalmw.Authorize("recipes", "write"))

		productRecipeGroup := authorizedGroup.Group("/products/:main_product_uuid/recipes", internalmw.Authorize("recipes", "read"))
		productRecipeGroup.GET("", recipeHandler.GetRecipesByMainProduct)

		// Outlet routes
		outletGroup := authorizedGroup.Group("/outlets", internalmw.Authorize("outlets", "read"))
		outletGroup.GET("", outletHandler.GetAllOutlets)
		outletGroup.GET("/:uuid", outletHandler.GetOutletByID)
		outletGroup.POST("", outletHandler.CreateOutlet, internalmw.Authorize("outlets", "write"), withValidation(&dtos.OutletCreateRequest{}, validators.ValidateCreateOutlet))
		outletGroup.PUT("/:uuid", outletHandler.UpdateOutlet, internalmw.Authorize("outlets", "write"), withValidation(&dtos.OutletUpdateRequest{}, validators.ValidateUpdateOutlet))
		outletGroup.DELETE("/:uuid", outletHandler.DeleteOutlet, internalmw.Authorize("outlets", "write"))

		// Stock routes
		stockGroup := authorizedGroup.Group("/outlets/:outlet_uuid/stocks", internalmw.Authorize("stocks", "read"))
		stockGroup.GET("", stockHandler.GetOutletStocks)
		stockGroup.GET("/:product_uuid", stockHandler.GetStockByOutletAndProduct)
		stockGroup.PUT("/:product_uuid", stockHandler.UpdateStock, internalmw.Authorize("stocks", "write"), withValidation(&dtos.UpdateStockRequest{}, validators.ValidateUpdateStock))

		authorizedGroup.PUT("/stocks", stockHandler.UpdateGlobalStock, internalmw.Authorize("stocks", "write"), withValidation(&dtos.GlobalStockUpdateRequest{}, validators.ValidateGlobalStockUpdate))

		// Order routes
		orderGroup := authorizedGroup.Group("/orders")
		orderGroup.POST("", orderHandler.CreateOrder, internalmw.Authorize("orders", "write"), withValidation(&dtos.CreateOrderRequest{}, validators.ValidateCreateOrder))
		orderGroup.GET("/:uuid", orderHandler.GetOrderByUuid, internalmw.Authorize("orders", "read"))

		outletOrdersGroup := authorizedGroup.Group("/outlets/:outlet_uuid/orders", internalmw.Authorize("orders", "read"))
		outletOrdersGroup.GET("", orderHandler.GetOrdersByOutlet)

		// Report routes
		reportGroup := authorizedGroup.Group("/reports", internalmw.Authorize("reports", "read"))
		reportGroup.GET("/outlets/:outlet_uuid/sales", reportHandler.GetSalesByOutletReport)
		reportGroup.GET("/products/:product_uuid/sales", reportHandler.GetSalesByProductReport)

		// Supplier routes
		supplierGroup := authorizedGroup.Group("/suppliers", internalmw.Authorize("suppliers", "read"))
		supplierGroup.GET("", supplierHandler.GetAllSuppliers)
		supplierGroup.GET("/:uuid", supplierHandler.GetSupplierByuuid)
		supplierGroup.POST("", supplierHandler.CreateSupplier, internalmw.Authorize("suppliers", "write"), withValidation(&dtos.CreateSupplierRequest{}, validators.ValidateCreateSupplier))
		supplierGroup.PUT("/:uuid", supplierHandler.UpdateSupplier, internalmw.Authorize("suppliers", "write"), withValidation(&dtos.UpdateSupplierRequest{}, validators.ValidateUpdateSupplier))
		supplierGroup.DELETE("/:uuid", supplierHandler.DeleteSupplier, internalmw.Authorize("suppliers", "write"))

		// Purchase Order routes
		poGroup := authorizedGroup.Group("/purchase-orders")
		poGroup.POST("", poHandler.CreatePurchaseOrder, internalmw.Authorize("purchase_orders", "write"), withValidation(&dtos.CreatePurchaseOrderRequest{}, validators.ValidateCreatePurchaseOrder))
		poGroup.PUT("/:uuid/receive", poHandler.ReceivePurchaseOrder, internalmw.Authorize("purchase_orders", "write"))
		poGroup.GET("/:uuid", poHandler.GetPurchaseOrderByUuid, internalmw.Authorize("purchase_orders", "read"))

		outletPoGroup := authorizedGroup.Group("/outlets/:outlet_uuid/purchase-orders", internalmw.Authorize("purchase_orders", "read"))
		outletPoGroup.GET("", poHandler.GetPurchaseOrdersByOutlet)
	}
}
