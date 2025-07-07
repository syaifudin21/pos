package resetdb

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/msyaifudin/pos/cmd/migrate"
	"github.com/msyaifudin/pos/cmd/seed"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/models"
)

func Run() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	database.InitDB()

	log.Println("Starting database reset...")

	// Drop all tables
	log.Println("Dropping all tables...")
	err := database.DB.Migrator().DropTable(
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
	)
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}
	log.Println("All tables dropped.")

	// Run migrations
	log.Println("Running migrations...")
	migrate.Run()

	// Run seeder
	log.Println("Running seeder...")
	seed.Run()

	log.Println("Database reset completed.")

	// Exit after reset
	os.Exit(0)
}
