package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/msyaifudin/pos/docs" // Import generated docs
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/handlers"
	internalmw "github.com/msyaifudin/pos/internal/middleware"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/pkg/casbin"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title POS API Documentation
// @version 1.0
// @description This is a sample POS API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	database.InitDB()

	// Auto-migrate models
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Outlet{},
		&models.Product{},
		&models.Recipe{},
		&models.Stock{},
		&models.Order{},
		&models.OrderItem{},
		&models.BaseModel{},
	) // BaseModel is included for its UUID type to be recognized by GORM
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database migration completed")

	// Initialize Casbin
	casbin.InitCasbin()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Inject DB into context for middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", database.DB)
			return next(c)
		}
	})

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Basic route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to POS API!")
	})

	// Initialize services and handlers
	authService := services.NewAuthService(database.DB)
	authHandler := handlers.NewAuthHandler(authService)

	// Auth routes
	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// User management routes (admin only)
	userAdminGroup := e.Group("/auth/users")
	userAdminGroup.Use(internalmw.Authorize("users", "manage")) // New Casbin rule for user management
	userAdminGroup.PUT("/:uuid/block", authHandler.BlockUser)
	userAdminGroup.PUT("/:uuid/unblock", authHandler.UnblockUser)

	// Initialize product and outlet services and handlers
	productService := services.NewProductService(database.DB)
	productHandler := handlers.NewProductHandler(productService)

	outletService := services.NewOutletService(database.DB)
	outletHandler := handlers.NewOutletHandler(outletService)

	// Product routes
	productGroup := e.Group("/products")
	productGroup.Use(internalmw.Authorize("products", "read"))
	productGroup.GET("/", productHandler.GetAllProducts)
	productGroup.GET("/:uuid", productHandler.GetProductByID)
	productGroup.Use(internalmw.Authorize("products", "write")) // For create, update, delete
	productGroup.POST("/", productHandler.CreateProduct)
	productGroup.PUT("/:uuid", productHandler.UpdateProduct)
	productGroup.DELETE("/:uuid", productHandler.DeleteProduct)

	// Initialize recipe service and handler
	recipeService := services.NewRecipeService(database.DB)
	recipeHandler := handlers.NewRecipeHandler(recipeService)

	// Recipe routes
	recipeGroup := e.Group("/recipes")
	recipeGroup.Use(internalmw.Authorize("recipes", "read"))
	recipeGroup.GET("/:uuid", recipeHandler.GetRecipeByUuid)
	recipeGroup.Use(internalmw.Authorize("recipes", "write")) // For create, update, delete
	recipeGroup.POST("/", recipeHandler.CreateRecipe)
	recipeGroup.PUT("/:uuid", recipeHandler.UpdateRecipe)
	recipeGroup.DELETE("/:uuid", recipeHandler.DeleteRecipe)

	// Recipes by main product
	productRecipeGroup := e.Group("/products/:main_product_uuid/recipes")
	productRecipeGroup.Use(internalmw.Authorize("recipes", "read"))
	productRecipeGroup.GET("/", recipeHandler.GetRecipesByMainProduct)

	// Outlet routes
	outletGroup := e.Group("/outlets")
	outletGroup.Use(internalmw.Authorize("outlets", "read"))
	outletGroup.GET("/", outletHandler.GetAllOutlets)
	outletGroup.GET("/:uuid", outletHandler.GetOutletByID)
	outletGroup.Use(internalmw.Authorize("outlets", "write")) // For create, update, delete
	outletGroup.POST("/", outletHandler.CreateOutlet)
	outletGroup.PUT("/:uuid", outletHandler.UpdateOutlet)
	outletGroup.DELETE("/:uuid", outletHandler.DeleteOutlet)

	// Initialize stock service and handler
	stockService := services.NewStockService(database.DB)
	stockHandler := handlers.NewStockHandler(stockService)

	// Stock routes
	stockGroup := e.Group("/outlets/:outlet_uuid/stocks")
	stockGroup.Use(internalmw.Authorize("stocks", "read"))
	stockGroup.GET("/", stockHandler.GetOutletStocks)
	stockGroup.GET("/:product_uuid", stockHandler.GetStockByOutletAndProduct)
	stockGroup.Use(internalmw.Authorize("stocks", "write")) // For update
	stockGroup.PUT("/:product_uuid", stockHandler.UpdateStock)

	// Global stock update route
	e.PUT("/stocks", stockHandler.UpdateGlobalStock, internalmw.Authorize("stocks", "write"))

	// Initialize order service and handler
	orderService := services.NewOrderService(database.DB, stockService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Order routes
	orderGroup := e.Group("/orders")
	orderGroup.Use(internalmw.Authorize("orders", "write")) // Cashier can create orders
	orderGroup.POST("/", orderHandler.CreateOrder)
	orderGroup.Use(internalmw.Authorize("orders", "read")) // Admin/Manager can read orders
	orderGroup.GET("/:uuid", orderHandler.GetOrderByUuid)

	// Orders by outlet
	outletOrdersGroup := e.Group("/outlets/:outlet_uuid/orders")
	outletOrdersGroup.Use(internalmw.Authorize("orders", "read"))
	outletOrdersGroup.GET("/", orderHandler.GetOrdersByOutlet)

	// Initialize report service and handler
	reportService := services.NewReportService(database.DB)
	reportHandler := handlers.NewReportHandler(reportService)

	// Report routes
	reportGroup := e.Group("/reports")
	reportGroup.Use(internalmw.Authorize("reports", "read"))
	reportGroup.GET("/outlets/:outlet_uuid/sales", reportHandler.GetSalesByOutletReport)
	reportGroup.GET("/products/:product_uuid/sales", reportHandler.GetSalesByProductReport)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
