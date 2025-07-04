package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/handlers"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/pkg/casbin"
	"github.com/msyaifudin/pos/pkg/utils"
	"gorm.io/gorm"
)

// Authorize checks if the current user has permission to access the requested resource.
func Authorize(obj, act string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Message: "Authorization header missing"})
			}

			tokenString := authHeader
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}

			claims, err := utils.ParseToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Message: "Invalid or expired token"})
			}

			// Get DB instance from context
			db, ok := c.Get("db").(*gorm.DB)
			if !ok || db == nil {
				return c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{Message: "Database connection not available"})
			}

			// Check if user is blocked
			user, err := services.NewAuthService(db).GetUserByuuid(claims.ID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{Message: "Failed to retrieve user information"})
			}
			if user.IsBlocked {
				return c.JSON(http.StatusForbidden, handlers.ErrorResponse{Message: "Your account has been blocked."})
			}

			// Get user role from claims
			userRole := claims.Role

			// Enforce Casbin policy
			ok, err = casbin.Enforcer.Enforce(userRole, obj, act)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{Message: "Authorization error"})
			}

			if !ok {
				return c.JSON(http.StatusForbidden, handlers.ErrorResponse{Message: "Forbidden: You don't have permission to access this resource"})
			}

			// Store user claims in context for later use
			c.Set("claims", claims)

			return next(c)
		}
	}
}
