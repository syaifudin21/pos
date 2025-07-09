package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/handlers"
	internalmw "github.com/msyaifudin/pos/internal/middleware"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services" // New import
	"github.com/msyaifudin/pos/internal/validators"
	"gorm.io/gorm" // New import
)

func RegisterApiRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize services and handlers for API routes
	userContextService := services.NewUserContextService(db)
	authService := services.NewAuthService(db)
	authHandler := handlers.NewAuthHandler(authService, userContextService)

	userPaymentService := services.NewUserPaymentService(db, userContextService)
	userPaymentHandler := handlers.NewUserPaymentHandler(userPaymentService, userContextService)

	productService := services.NewProductService(db, userContextService)
	productHandler := handlers.NewProductHandler(productService, userContextService)

	recipeService := services.NewRecipeService(db, userContextService)
	recipeHandler := handlers.NewRecipeHandler(recipeService, userContextService)

	outletService := services.NewOutletService(db, userContextService)
	outletHandler := handlers.NewOutletHandler(outletService, userContextService)

	stockMovementService := services.NewStockMovementService(db)
	stockService := services.NewStockService(db, userContextService, stockMovementService)
	stockHandler := handlers.NewStockHandler(stockService, userContextService)

	ipaymuService := services.NewIpaymuService(db, userContextService) // Assuming this is needed for orderService
	orderService := services.NewOrderService(db, stockService, ipaymuService, userContextService)
	orderHandler := handlers.NewOrderHandler(orderService, userContextService)

	reportService := services.NewReportService(db)
	reportHandler := handlers.NewReportHandler(reportService, userContextService)

	supplierService := services.NewSupplierService(db, userContextService)
	supplierHandler := handlers.NewSupplierHandler(supplierService, userContextService)

	poService := services.NewPurchaseOrderService(db, stockService, userContextService)
	poHandler := handlers.NewPurchaseOrderHandler(poService, userContextService)

	authorizedGroup := e.Group("")
	{
		// User management routes (owner only)
		userAdminGroup := authorizedGroup.Group("/users", internalmw.Authorize("users", "manage"))
		userAdminGroup.GET("", authHandler.GetAllUsers)
		userAdminGroup.POST("", authHandler.Register, WithValidation(&dtos.RegisterRequest{}, validators.ValidateRegisterRequest))
		userAdminGroup.PUT("/:uuid", authHandler.UpdateUser, WithValidation(&dtos.UpdateUserRequest{}, validators.ValidateUpdateUserRequest))
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
		productGroup.POST("", productHandler.CreateProduct, internalmw.Authorize("products", "write"), WithValidation(&dtos.ProductCreateRequest{}, validators.ValidateCreateProduct))
		productGroup.PUT("/:uuid", productHandler.UpdateProduct, internalmw.Authorize("products", "write"), WithValidation(&dtos.ProductUpdateRequest{}, validators.ValidateUpdateProduct))
		productGroup.DELETE("/:uuid", productHandler.DeleteProduct, internalmw.Authorize("products", "write"))

		outletProductGroup := authorizedGroup.Group("/outlets/:outlet_uuid/products", internalmw.Authorize("products", "read"))
		outletProductGroup.GET("", productHandler.GetProductsByOutlet)

		// Recipe routes
		recipeGroup := authorizedGroup.Group("/recipes", internalmw.Authorize("recipes", "read"))
		recipeGroup.GET("/:uuid", recipeHandler.GetRecipeByUuid)
		recipeGroup.POST("", recipeHandler.CreateRecipe, internalmw.Authorize("recipes", "write"), WithValidation(&dtos.CreateRecipeRequest{}, validators.ValidateCreateRecipe))
		recipeGroup.PUT("/:uuid", recipeHandler.UpdateRecipe, internalmw.Authorize("recipes", "write"), WithValidation(&dtos.UpdateRecipeRequest{}, validators.ValidateUpdateRecipe))
		recipeGroup.DELETE("/:uuid", recipeHandler.DeleteRecipe, internalmw.Authorize("recipes", "write"))

		productRecipeGroup := authorizedGroup.Group("/products/:main_product_uuid/recipes", internalmw.Authorize("recipes", "read"))
		productRecipeGroup.GET("", recipeHandler.GetRecipesByMainProduct)

		// Outlet routes
		outletGroup := authorizedGroup.Group("/outlets", internalmw.Authorize("outlets", "read"))
		outletGroup.GET("", outletHandler.GetAllOutlets)
		outletGroup.GET("/:uuid", outletHandler.GetOutletByID)
		outletGroup.POST("", outletHandler.CreateOutlet, internalmw.Authorize("outlets", "write"), WithValidation(&dtos.OutletCreateRequest{}, validators.ValidateCreateOutlet))
		outletGroup.PUT("/:uuid", outletHandler.UpdateOutlet, internalmw.Authorize("outlets", "write"), WithValidation(&dtos.OutletUpdateRequest{}, validators.ValidateUpdateOutlet))
		outletGroup.DELETE("/:uuid", outletHandler.DeleteOutlet, internalmw.Authorize("outlets", "write"))

		// Stock routes
		stockGroup := authorizedGroup.Group("/outlets/:outlet_uuid/stocks", internalmw.Authorize("stocks", "read"))
		stockGroup.GET("", stockHandler.GetOutletStocks)
		stockGroup.GET("/:product_uuid", stockHandler.GetStockByOutletAndProduct)
		stockGroup.PUT("/:product_uuid", stockHandler.UpdateStock, internalmw.Authorize("stocks", "write"), WithValidation(&dtos.UpdateStockRequest{}, validators.ValidateUpdateStock))

		authorizedGroup.PUT("/stocks", stockHandler.UpdateGlobalStock, internalmw.Authorize("stocks", "write"), WithValidation(&dtos.GlobalStockUpdateRequest{}, validators.ValidateGlobalStockUpdate))

		// Order routes
		orderGroup := authorizedGroup.Group("/orders")
		orderGroup.POST("", orderHandler.CreateOrder, internalmw.Authorize("orders", "write"), WithValidation(&dtos.CreateOrderRequest{}, validators.ValidateCreateOrder))
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
		supplierGroup.POST("", supplierHandler.CreateSupplier, internalmw.Authorize("suppliers", "write"), WithValidation(&dtos.CreateSupplierRequest{}, validators.ValidateCreateSupplier))
		supplierGroup.PUT("/:uuid", supplierHandler.UpdateSupplier, internalmw.Authorize("suppliers", "write"), WithValidation(&dtos.UpdateSupplierRequest{}, validators.ValidateUpdateSupplier))
		supplierGroup.DELETE("/:uuid", supplierHandler.DeleteSupplier, internalmw.Authorize("suppliers", "write"))

		// Purchase Order routes
		poGroup := authorizedGroup.Group("/purchase-orders")
		poGroup.POST("", poHandler.CreatePurchaseOrder, internalmw.Authorize("purchase_orders", "write"), WithValidation(&dtos.CreatePurchaseOrderRequest{}, validators.ValidateCreatePurchaseOrder))
		poGroup.PUT("/:uuid/receive", poHandler.ReceivePurchaseOrder, internalmw.Authorize("purchase_orders", "write"))
		poGroup.GET("/:uuid", poHandler.GetPurchaseOrderByUuid, internalmw.Authorize("purchase_orders", "read"))

		outletPoGroup := authorizedGroup.Group("/outlets/:outlet_uuid/purchase-orders", internalmw.Authorize("purchase_orders", "read"))
		outletPoGroup.GET("", poHandler.GetPurchaseOrdersByOutlet)
	}
}
