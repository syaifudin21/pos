package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/database"
)

func RegisterRoutes(e *echo.Echo) {
	// Pass database connection to route registration functions
	RegisterPublicRoutes(e, database.DB)
	RegisterAuthRoutes(e, database.DB)
	RegisterAccountRoutes(e, database.DB)
	RegisterApiRoutes(e, database.DB)
}
