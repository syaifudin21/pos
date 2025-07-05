package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/pkg/utils"
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

	// Only allow 'manager' and 'cashier' roles for this endpoint
	if req.Role != "manager" && req.Role != "cashier" {
		return JSONError(c, http.StatusBadRequest, "invalid_role_for_this_registration_endpoint")
	}

	// Get the ID of the currently logged-in admin from the JWT claims
	claims := c.Get("user").(*jwt.Token).Claims.(*utils.Claims)
	creatorID := claims.ID

	user, err := h.AuthService.RegisterUser(req.Username, req.Password, req.Role, req.OutletID, &creatorID, nil, nil)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "user_registered_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role})
}

func (h *AuthHandler) RegisterAdmin(c echo.Context) error {
	req := new(dtos.RegisterAdminRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	if req.Username == "" || req.Password == "" || req.Email == "" || req.PhoneNumber == "" {
		return JSONError(c, http.StatusBadRequest, "username_password_email_phone_required")
	}

	// Only allow 'admin' role for this endpoint
	// For initial admin registration, this check might be skipped or handled differently
	// For subsequent admin registrations by an existing admin, you might want to check the creator's role
	// For simplicity, assuming this is for initial admin setup or by a super-admin

	// No creatorID for the first admin, or if registered by a super-admin outside the system
	// For now, let's assume no creatorID for admin registration via this endpoint
	user, err := h.AuthService.RegisterUser(req.Username, req.Password, "admin", nil, nil, &req.Email, &req.PhoneNumber)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "admin_registered_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role})
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

	return JSONSuccess(c, http.StatusOK, "login_successful", dtos.LoginResponse{Token: token, User: dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role}})
}

func (h *AuthHandler) BlockUser(c echo.Context) error {
	userUuidParam := c.Param("uuid")
	userUuid, err := uuid.Parse(userUuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_user_uuid_format")
	}

	user, err := h.AuthService.GetUserByuuid(userUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	if err := h.AuthService.BlockUser(user.ID); err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "user_blocked_successfully", nil)
}

func (h *AuthHandler) UnblockUser(c echo.Context) error {
	userUuidParam := c.Param("uuid")
	userUuid, err := uuid.Parse(userUuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_user_uuid_format")
	}

	user, err := h.AuthService.GetUserByuuid(userUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	if err := h.AuthService.UnblockUser(user.ID); err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "user_unblocked_successfully", nil)
}

func (h *AuthHandler) GetAllUsers(c echo.Context) error {
	claims := c.Get("user").(*jwt.Token).Claims.(*utils.Claims)
	adminID := claims.ID

	users, err := h.AuthService.GetAllUsers(adminID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	var userResponses []dtos.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role})
	}
	return JSONSuccess(c, http.StatusOK, "users_retrieved_successfully", userResponses)
}

func (h *AuthHandler) UpdateUser(c echo.Context) error {
	userUuidParam := c.Param("uuid")
	userUuid, err := uuid.Parse(userUuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_user_uuid_format")
	}

	user, err := h.AuthService.GetUserByuuid(userUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	req := new(dtos.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	updatedUser, err := h.AuthService.UpdateUser(user.ID, req)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "user_updated_successfully", dtos.UserResponse{ID: updatedUser.ID, Uuid: updatedUser.Uuid, Username: updatedUser.Username, Role: updatedUser.Role})
}
