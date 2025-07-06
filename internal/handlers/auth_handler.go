package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/validators"
	"github.com/msyaifudin/pos/pkg/utils"
)

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

	lang := c.Get("lang").(string)
	if messages := validators.ValidateRegisterRequest(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	// Get the ID of the currently logged-in admin from the JWT claims
	claims := c.Get("user").(*jwt.Token).Claims.(*utils.Claims)
	creatorID := claims.ID

	user, err := h.AuthService.RegisterUser(req.Password, req.Role, req.OutletID, &creatorID, nil, nil)
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

	lang := c.Get("lang").(string)
	if messages := validators.ValidateRegisterAdminRequest(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	// Only allow 'admin' role for this endpoint
	// For initial admin registration, this check might be skipped or handled differently
	// For subsequent admin registrations by an existing admin, you might want to check the creator's role
	// For simplicity, assuming this is for initial admin setup or by a super-admin

	// No creatorID for the first admin, or if registered by a super-admin outside the system
	// For now, let's assume no creatorID for admin registration via this endpoint
	user, err := h.AuthService.RegisterUser(req.Password, "admin", nil, nil, &req.Email, &req.PhoneNumber)
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

	lang := c.Get("lang").(string)
	if messages := validators.ValidateLoginRequest(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
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
	req := new(dtos.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateUpdateUserRequest(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	_, err = h.AuthService.GetUserByuuid(userUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "user_updated_successfully", nil)
}

func (h *AuthHandler) DeleteUser(c echo.Context) error {
	userUuidParam := c.Param("uuid")
	userUuid, err := uuid.Parse(userUuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_user_uuid_format")
	}

	user, err := h.AuthService.GetUserByuuid(userUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	if err := h.AuthService.DeleteUser(user.ID); err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "user_deleted_successfully", nil)
}

func (h *AuthHandler) VerifyOTP(c echo.Context) error {
	req := new(dtos.VerifyOTPRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateVerifyOTPRequest(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user, err := h.AuthService.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "otp_verified_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role})
}
