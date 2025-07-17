package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/pkg/localization"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// JSONSuccess sends a success JSON response.
func JSONSuccess(c echo.Context, statusCode int, messageKey string, data interface{}) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizedMessage := localization.GetLocalizedMessage(messageKey, lang)
	return c.JSON(statusCode, SuccessResponse{Message: localizedMessage, Data: data})
}

// JSONError sends an error JSON response or redirects.
func JSONError(c echo.Context, statusCode int, messageKey interface{}, redirectURL ...string) error {
	lang := c.Request().Header.Get("Accept-Language")

	var localizedMessage string
	var details interface{}

	switch msg := messageKey.(type) {
	case string:
		localizedMessage = localization.GetLocalizedMessage(msg, lang)
	case validator.ValidationErrors:
		localizedMessage = localization.GetLocalizedMessage("validation_error", lang)
		errors := make(map[string]string)
		for _, err := range msg {
			errors[err.Field()] = localization.GetLocalizedMessage(err.Tag(), lang)
		}
		details = errors
	default:
		localizedMessage = localization.GetLocalizedMessage("unknown_error", lang)
	}

	if len(redirectURL) > 0 && redirectURL[0] != "" {
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL[0])
	}

	return c.JSON(statusCode, ErrorResponse{Message: localizedMessage, Details: details})
}

// MapErrorToStatusCode maps common error messages to HTTP status codes.
func MapErrorToStatusCode(err error) int {
	switch err.Error() {
	case "user not found", "outlet not found", "product not found", "supplier not found", "recipe not found", "stock not found", "order not found", "purchase order not found":
		return http.StatusNotFound
	case "invalid credentials", "unauthorized", "user not verified":
		return http.StatusUnauthorized
	case "username already exists", "invalid input", "validation error":
		return http.StatusBadRequest
	case "forbidden":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
