package database

import (
	"fmt"
	"log"
	"os"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta",
		getEnvOrDefault("DB_USER", "msyaifudin"),
		getEnvOrDefault("DB_PASSWORD", ""),
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_PORT", "5432"),
		getEnvOrDefault("DB_NAME", "yespos2"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	if err := db.Use(&UpdateByCallback{}); err != nil {
		log.Fatalf("Gagal register plugin UpdateByCallback: %v", err)
	}

	// AutoMigrate all models, including CasbinRule
	if err := db.AutoMigrate(&models.IpaymuLog{}, &gormadapter.CasbinRule{}, &models.Order{}, &models.OrderPayment{}); err != nil {
		log.Fatalf("Gagal migrasi tabel: %v", err)
	}

	DB = db

	log.Println("Database connection established")
}

func getEnvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
