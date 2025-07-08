package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/msyaifudin/pos/cmd/migrate"
	"github.com/msyaifudin/pos/internal/database"
	internalmw "github.com/msyaifudin/pos/internal/middleware"
	"github.com/msyaifudin/pos/internal/redis"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/routes"
	"github.com/msyaifudin/pos/pkg/casbin"
)

func Run() {
	// Perform database migrations
	migrate.PerformMigration()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Redis
	redis.InitRedis()

	// Initialize Casbin
	casbin.InitCasbin()

	// Initialize Email Queue and Worker
	services.InitEmailQueue()
	services.StartEmailWorker()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(30),
	)))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return internalmw.LanguageMiddleware(next)
	})
	// Inject DB into context for middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", database.DB)
			return next(c)
		}
	})

	// Register all routes
	routes.RegisterRoutes(e)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
