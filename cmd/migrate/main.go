package migrate

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/models"
)

func PerformMigration() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	database.InitDB()

	log.Println("Starting database migration...")

	// Auto-migrate models
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Outlet{},
		&models.Product{},
		&models.Recipe{},
		&models.Stock{},
		&models.Order{},
		&models.OrderItem{},
		&models.Supplier{},
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.OTP{},
		&models.PaymentMethod{},
		&models.IpaymuLog{},
		&models.UserPayment{},
		&models.UserIpaymu{},
		&models.UserTsm{},
		&models.StockMovement{},
		&models.ProductVariant{},
		&models.OrderItemAddOn{},
		&models.ProductAddOn{},
		&models.TsmLog{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database migration completed.")
}

func Run() {
	PerformMigration()
	// Exit after migration
	os.Exit(0)
}
