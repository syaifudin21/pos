package middleware

import (
	"github.com/labstack/echo/v4"
)

// LanguageMiddleware extracts the Accept-Language header and sets it in the Echo context.
// If the header is not present, it defaults to "en".
func LanguageMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lang := c.Request().Header.Get("Accept-Language")
		if lang == "" {
			lang = "en" // Default language if not provided
		}
		c.Set("lang", lang)
		return next(c)
	}
}
