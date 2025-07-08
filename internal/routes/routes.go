package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/handlers"
	internalmw "github.com/msyaifudin/pos/internal/middleware"
	"github.com/msyaifudin/pos/internal/services"
)

func RegisterRoutes(e *echo.Echo) {
	// Initialize iPaymu service
	ipaymuService := services.NewIpaymuService(database.DB)

	// Initialize services and handlers
	userContextService := services.NewUserContextService(database.DB)
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
	ipaymuHandler := handlers.NewIpaymuHandler(ipaymuService)
	userPaymentService := services.NewUserPaymentService(database.DB, userContextService)
	userPaymentHandler := handlers.NewUserPaymentHandler(userPaymentService, userContextService)
	tsmService := services.NewTsmService(database.DB, userContextService)
	tsmHandler := handlers.NewTsmHandler(tsmService, userContextService)

	// Public routes (no specific middleware)
	e.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to POS API!")
	})
	e.POST("/api/payment/ipaymu/notify", ipaymuHandler.IpaymuNotify)

	// Auth routes (public, but grouped)
	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.RegisterOwner)
	authGroup.POST("/verify-otp", authHandler.VerifyOTP)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/forgot-password", authHandler.ForgotPassword)
	authGroup.POST("/reset-password", authHandler.ResetPassword)
	authGroup.POST("/resend-verification-email", authHandler.ResendVerificationEmail)
	authGroup.GET("/google/login", googleOAuthHandler.GoogleLogin)
	authGroup.GET("/google/callback", googleOAuthHandler.GoogleCallback)

	// Routes requiring internalmw.SelfAuthorize()
	selfAuthGroup := e.Group("")
	selfAuthGroup.Use(internalmw.SelfAuthorize())
	{
		accountGroup := selfAuthGroup.Group("/account")
		accountGroup.GET("/profile", authHandler.GetProfile)
		accountGroup.PUT("/password", authHandler.UpdatePassword)
		accountGroup.POST("/email/otp", authHandler.SendOTPForEmailUpdate)
		accountGroup.PUT("/email", authHandler.UpdateEmail)

		paymentGroup := selfAuthGroup.Group("/ipaymu")
		paymentGroup.POST("/register", ipaymuHandler.RegisterIpaymu)
		paymentGroup.POST("/direct-payment", ipaymuHandler.CreateDirectPayment)

		tsmGroup := selfAuthGroup.Group("/tsm")
		tsmGroup.POST("/register", tsmHandler.RegisterTsm)
	}

	// Routes requiring internalmw.Authorize()
	authorizedGroup := e.Group("")
	{
		// User management routes (owner only)
		userAdminGroup := authorizedGroup.Group("/users")
		userAdminGroup.Use(internalmw.Authorize("users", "manage"))
		userAdminGroup.GET("", authHandler.GetAllUsers)
		userAdminGroup.POST("", authHandler.Register)
		userAdminGroup.PUT("/:uuid", authHandler.UpdateUser)
		userAdminGroup.PUT("/:uuid/block", authHandler.BlockUser)
		userAdminGroup.PUT("/:uuid/unblock", authHandler.UnblockUser)
		userAdminGroup.DELETE("/:uuid", authHandler.DeleteUser)

		// User Payment routes (owner only)
		userPaymentGroup := authorizedGroup.Group("/account/payment-methods")
		userPaymentGroup.POST("/activate", userPaymentHandler.ActivateUserPayment, internalmw.Authorize("user_payments", "activate"))
		userPaymentGroup.POST("/deactivate", userPaymentHandler.DeactivateUserPayment, internalmw.Authorize("user_payments", "deactivate"))
		userPaymentGroup.GET("", userPaymentHandler.ListUserPayments, internalmw.Authorize("user_payments", "read"))

		// Product routes
		productGroup := authorizedGroup.Group("/products")
		productGroup.Use(internalmw.Authorize("products", "read"))
		productGroup.GET("", productHandler.GetAllProducts)
		productGroup.GET("/:uuid", productHandler.GetProductByID)
		productGroup.Use(internalmw.Authorize("products", "write"))
		productGroup.POST("", productHandler.CreateProduct)
		productGroup.PUT("/:uuid", productHandler.UpdateProduct)
		productGroup.DELETE("/:uuid", productHandler.DeleteProduct)

		outletProductGroup := authorizedGroup.Group("/outlets/:outlet_uuid/products")
		outletProductGroup.Use(internalmw.Authorize("products", "read"))
		outletProductGroup.GET("", productHandler.GetProductsByOutlet)

		// Recipe routes
		recipeGroup := authorizedGroup.Group("/recipes")
		recipeGroup.Use(internalmw.Authorize("recipes", "read"))
		recipeGroup.GET("/:uuid", recipeHandler.GetRecipeByUuid)
		recipeGroup.Use(internalmw.Authorize("recipes", "write"))
		recipeGroup.POST("", recipeHandler.CreateRecipe)
		recipeGroup.PUT("/:uuid", recipeHandler.UpdateRecipe)
		recipeGroup.DELETE("/:uuid", recipeHandler.DeleteRecipe)

		productRecipeGroup := authorizedGroup.Group("/products/:main_product_uuid/recipes")
		productRecipeGroup.Use(internalmw.Authorize("recipes", "read"))
		productRecipeGroup.GET("", recipeHandler.GetRecipesByMainProduct)

		// Outlet routes
		outletGroup := authorizedGroup.Group("/outlets")
		outletGroup.Use(internalmw.Authorize("outlets", "read"))
		outletGroup.GET("", outletHandler.GetAllOutlets)
		outletGroup.GET("/:uuid", outletHandler.GetOutletByID)
		outletGroup.Use(internalmw.Authorize("outlets", "write"))
		outletGroup.POST("", outletHandler.CreateOutlet)
		outletGroup.PUT("/:uuid", outletHandler.UpdateOutlet)
		outletGroup.DELETE("/:uuid", outletHandler.DeleteOutlet)

		// Stock routes
		stockGroup := authorizedGroup.Group("/outlets/:outlet_uuid/stocks")
		stockGroup.Use(internalmw.Authorize("stocks", "read"))
		stockGroup.GET("", stockHandler.GetOutletStocks)
		stockGroup.GET("/:product_uuid", stockHandler.GetStockByOutletAndProduct)
		stockGroup.Use(internalmw.Authorize("stocks", "write"))
		stockGroup.PUT("/:product_uuid", stockHandler.UpdateStock)

		authorizedGroup.PUT("/stocks", stockHandler.UpdateGlobalStock, internalmw.Authorize("stocks", "write"))

		// Order routes
		orderGroup := authorizedGroup.Group("/orders")
		orderGroup.Use(internalmw.Authorize("orders", "write"))
		orderGroup.POST("", orderHandler.CreateOrder)
		orderGroup.Use(internalmw.Authorize("orders", "read"))
		orderGroup.GET("/:uuid", orderHandler.GetOrderByUuid)

		outletOrdersGroup := authorizedGroup.Group("/outlets/:outlet_uuid/orders")
		outletOrdersGroup.Use(internalmw.Authorize("orders", "read"))
		outletOrdersGroup.GET("", orderHandler.GetOrdersByOutlet)

		// Report routes
		reportGroup := authorizedGroup.Group("/reports")
		reportGroup.Use(internalmw.Authorize("reports", "read"))
		reportGroup.GET("/outlets/:outlet_uuid/sales", reportHandler.GetSalesByOutletReport)
		reportGroup.GET("/products/:product_uuid/sales", reportHandler.GetSalesByProductReport)

		// Supplier routes
		supplierGroup := authorizedGroup.Group("/suppliers")
		supplierGroup.Use(internalmw.Authorize("suppliers", "read"))
		supplierGroup.GET("", supplierHandler.GetAllSuppliers)
		supplierGroup.GET("/:uuid", supplierHandler.GetSupplierByuuid)
		supplierGroup.Use(internalmw.Authorize("suppliers", "write"))
		supplierGroup.POST("", supplierHandler.CreateSupplier)
		supplierGroup.PUT("/:uuid", supplierHandler.UpdateSupplier)
		supplierGroup.DELETE("/:uuid", supplierHandler.DeleteSupplier)

		// Purchase Order routes
		poGroup := authorizedGroup.Group("/purchase-orders")
		poGroup.Use(internalmw.Authorize("purchase_orders", "write"))
		poGroup.POST("", poHandler.CreatePurchaseOrder)
		poGroup.PUT("/:uuid/receive", poHandler.ReceivePurchaseOrder)
		poGroup.Use(internalmw.Authorize("purchase_orders", "read"))
		poGroup.GET("/:uuid", poHandler.GetPurchaseOrderByUuid)

		outletPoGroup := authorizedGroup.Group("/outlets/:outlet_uuid/purchase-orders")
		outletPoGroup.Use(internalmw.Authorize("purchase_orders", "read"))
		outletPoGroup.GET("", poHandler.GetPurchaseOrdersByOutlet)
	}
}
