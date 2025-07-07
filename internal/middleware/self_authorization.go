package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/handlers"
	"github.com/msyaifudin/pos/pkg/utils"
)

// SelfAuthorize ensures that the authenticated user is accessing/modifying their own data.
// This middleware assumes that the user ID is available in the JWT claims.
func SelfAuthorize() echo.MiddlewareFunc {
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

			token, err := utils.ParseToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Message: "Invalid or expired token"})
			}

			claims := token.Claims.(*utils.Claims)

			// Store user claims in context for later use
			c.Set("user", token)
			c.Set("userID", claims.ID) // Set the user ID in context for easy access

			return next(c)
		}
	}
}
