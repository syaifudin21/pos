package seed

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/pkg/utils"
	"gorm.io/gorm"
)

func Run() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	database.InitDB()

	log.Println("Starting database seeding...")

	// Seed Payment Methods
	seedPaymentMethods(database.DB)

	// Seed Super Admin User
	seedSuperAdmin(database.DB)

	log.Println("Database seeding completed.")
}

func seedPaymentMethods(db *gorm.DB) {
	paymentMethods := []models.PaymentMethod{
		{Name: "Cash", Type: "cash", IsActive: true, PaymentMethod: "cash", PaymentChannel: "manual"},
		{Name: "Bank Transfer", Type: "bank_transfer", IsActive: true, PaymentMethod: "va", PaymentChannel: "mandiri"},
		{Name: "Credit Card", Type: "credit_card", IsActive: true, PaymentMethod: "edc", PaymentChannel: "linkpayment"},
		{Name: "QRIS", Type: "qris", IsActive: true, PaymentMethod: "qris", PaymentChannel: "qris"},
	}

	for _, pm := range paymentMethods {
		var existingPM models.PaymentMethod
		if err := db.Where("name = ?", pm.Name).First(&existingPM).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&pm).Error; err != nil {
					log.Printf("Failed to seed payment method %s: %v", pm.Name, err)
				} else {
					log.Printf("Seeded payment method: %s", pm.Name)
				}
			} else {
				log.Printf("Error checking existing payment method %s: %v", pm.Name, err)
			}
		} else {
			log.Printf("Payment method %s already exists, skipping.", pm.Name)
		}
	}
}

func seedSuperAdmin(db *gorm.DB) {
	email := "super@super.com"
	password := "S@ndi1234"
	name := "Super Admin"

	var existingUser models.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hashedPassword, err := utils.HashPassword(password)
			if err != nil {
				log.Fatalf("Failed to hash super admin password: %v", err)
			}

			superAdmin := models.User{
				Name:            name,
				Email:           email,
				Password:        hashedPassword,
				Role:            "admin",
				EmailVerifiedAt: func() *time.Time { t := time.Now(); return &t }(), // Mark as verified
			}

			if err := db.Create(&superAdmin).Error; err != nil {
				log.Fatalf("Failed to seed super admin user: %v", err)
			}
			log.Println("Seeded super admin user:", email)

			// Find the 'Cash' payment method
			var cashPaymentMethod models.PaymentMethod
			if err := db.Where("name = ?", "Cash").First(&cashPaymentMethod).Error; err != nil {
				log.Printf("Error finding 'Cash' payment method: %v", err)
			} else {
				// Create a UserPayment entry for the super admin with 'Cash' payment method
				userPayment := models.UserPayment{
					UserID:          superAdmin.ID,
					PaymentMethodID: cashPaymentMethod.ID,
					IsActive:        true,
				}
				if err := db.Create(&userPayment).Error; err != nil {
					log.Printf("Error creating default UserPayment for super admin: %v", err)
				} else {
					log.Println("Created default UserPayment for super admin with 'Cash' payment method.")
				}
			}

		} else {
			log.Fatalf("Error checking existing super admin user: %v", err)
		}
	} else {
		log.Println("Super admin user already exists, skipping.")
	}
}
