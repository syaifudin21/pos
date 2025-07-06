package handlers

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/services"
)

type GoogleOAuthHandler struct {
	GoogleOAuthService *services.GoogleOAuthService
}

func NewGoogleOAuthHandler(googleOAuthService *services.GoogleOAuthService) *GoogleOAuthHandler {
	return &GoogleOAuthHandler{GoogleOAuthService: googleOAuthService}
}

func (h *GoogleOAuthHandler) GoogleLogin(c echo.Context) error {
	url := h.GoogleOAuthService.GetGoogleLoginURL()
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *GoogleOAuthHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return JSONError(c, http.StatusBadRequest, "authorization_code_missing", os.Getenv("HOST"))
	}

	token, user, err := h.GoogleOAuthService.HandleGoogleCallback(code)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error(), os.Getenv("HOST"))
	}

	// Redirect or return JWT token
	return JSONSuccess(c, http.StatusOK, "google_login_successful", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
