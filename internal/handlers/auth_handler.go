package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

func isValidRole(role string) bool {
	for _, r := range models.AllowedUserRoles {
		if r == role {
			return true
		}
	}
	return false
}

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(dtos.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	if req.Username == "" || req.Password == "" || req.Role == "" {
		return JSONError(c, http.StatusBadRequest, "username_password_required")
	}

	if !isValidRole(req.Role) {
		return JSONError(c, http.StatusBadRequest, "invalid_role_specified")
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Password, req.Role, req.OutletID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "user_registered_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role, OutletID: user.OutletID})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(dtos.LoginRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	if req.Username == "" || req.Password == "" {
		return JSONError(c, http.StatusBadRequest, "username_password_required")
	}

	token, user, err := h.AuthService.LoginUser(req.Username, req.Password)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "login_successful", dtos.LoginResponse{Token: token, User: dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role, OutletID: user.OutletID}})
}

func (h *AuthHandler) BlockUser(c echo.Context) error {
	useruuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_user_uuid_format")
	}

	if err := h.AuthService.BlockUser(useruuid); err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "user_blocked_successfully", nil)
}

func (h *AuthHandler) UnblockUser(c echo.Context) error {
	useruuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_user_uuid_format")
	}

	if err := h.AuthService.UnblockUser(useruuid); err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "user_unblocked_successfully", nil)
}
