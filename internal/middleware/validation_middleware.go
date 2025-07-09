package middleware

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/handlers"
	"github.com/msyaifudin/pos/pkg/localization"
)

// ValidationMiddleware is a generic middleware for request validation.
// It takes an empty DTO instance and a validator function.
func ValidationMiddleware(dtoType interface{}, validatorFunc func(interface{}) []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create a new instance of the DTO type using reflection
			val := reflect.New(reflect.TypeOf(dtoType).Elem())
			req := val.Interface()

			if err := c.Bind(req); err != nil {
				// Check if it's a binding error (e.g., JSON parsing, type mismatch)
				if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusBadRequest {
					return handlers.JSONError(c, http.StatusBadRequest, "Invalid JSON format or data type mismatch.")
				}
				return handlers.JSONError(c, http.StatusBadRequest, "invalid_input")
			}

			lang := c.Get("lang").(string)
			if messageKeys := validatorFunc(req); messageKeys != nil {
				localizedMessages := make([]string, 0, len(messageKeys))
				for _, key := range messageKeys {
					localizedMessages = append(localizedMessages, localization.GetLocalizedMessage(key, lang))
				}
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": localizedMessages,
				})
			}

			// Store the validated request in context
			c.Set("validated_data", req)

			return next(c)
		}
	}
}
