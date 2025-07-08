package seed

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/pkg/casbin"
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

	// Seed Casbin Policies from CSV to DB
	seedCasbinPolicies()

	// Seed Super Admin User
	seedSuperAdmin(database.DB)

	log.Println("Database seeding completed.")
}

func seedPaymentMethods(db *gorm.DB) {
	paymentMethods := []models.PaymentMethod{
		{Issuer: "default", Name: "Cash", Type: "cash", IsActive: true, PaymentMethod: "cash", PaymentChannel: "manual"},
		{Issuer: "iPaymu", Name: "Bank Transfer", Type: "bank_transfer", IsActive: true, PaymentMethod: "va", PaymentChannel: "mandiri"},
		{Issuer: "TSM", Name: "Credit Card", Type: "credit_card", IsActive: true, PaymentMethod: "edc", PaymentChannel: "linkpayment"},
		{Issuer: "iPaymu", Name: "QRIS", Type: "qris", IsActive: true, PaymentMethod: "qris", PaymentChannel: "qris"},
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

func seedCasbinPolicies() {
	// Initialize Casbin enforcer with GORM adapter
	casbin.InitCasbin()

	// Clear existing policies in DB to avoid duplicates during seeding, if enabled by environment variable
	if os.Getenv("CASBIN_CLEAR_POLICY_ON_SEED") == "true" {
		log.Println("Clearing existing Casbin policies from DB...")
		casbin.Enforcer.ClearPolicy()
	} else {
		log.Println("Skipping clearing Casbin policies. To clear, set CASBIN_CLEAR_POLICY_ON_SEED=true")
	}

	file, err := os.Open("pkg/casbin/policy-default.csv")
	if err != nil {
		log.Fatalf("Failed to open policy.csv: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		parts := strings.Split(line, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		policyType := parts[0]
		policyArgs := parts[1:]

		switch policyType {
		case "p":
			// 'p' policies expect 3 arguments: sub, obj, act
			if len(policyArgs) != 3 {
				log.Printf("Skipping invalid 'p' policy line (expected 3 arguments, got %d): %s", len(policyArgs), line)
				continue
			}
			args := make([]interface{}, len(policyArgs))
			for i, v := range policyArgs {
				args[i] = v
			}
			if _, err := casbin.Enforcer.AddPolicy(args...); err != nil {
				log.Fatalf("Failed to add policy %v: %v", policyArgs, err)
			}
		case "g":
			// 'g' policies expect 2 arguments: user, role
			if len(policyArgs) != 2 {
				log.Printf("Skipping invalid 'g' policy line (expected 2 arguments, got %d): %s", len(policyArgs), line)
				continue
			}
			args := make([]interface{}, len(policyArgs))
			for i, v := range policyArgs {
				args[i] = v
			}
			if _, err := casbin.Enforcer.AddGroupingPolicy(args...); err != nil {
				log.Fatalf("Failed to add grouping policy %v: %v", policyArgs, err)
			}
		default:
			log.Printf("Unknown policy type: %s in line: %s", policyType, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading policy.csv: %v", err)
	}

	log.Println("Casbin policies seeded from policy.csv to database.")

	if err := casbin.UpdateCasbinPolicy(); err != nil {
		log.Fatalf("Failed to update Casbin policy via Redis watcher: %v", err)
	}
}
