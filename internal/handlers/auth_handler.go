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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, password, role, and optional outlet ID.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dtos.RegisterRequest true "User registration details"
// @Success 201 {object} SuccessResponse{data=dtos.UserResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
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

// Login godoc
// @Summary Log in a user
// @Description Authenticate user and return a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body dtos.LoginRequest true "User login credentials"
// @Success 200 {object} SuccessResponse{data=dtos.LoginResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
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

// BlockUser godoc
// @Summary Block a user
// @Description Block a user by their Uuid. Only admin can perform this action.
// @Tags Auth
// @Accept json
// @Produce json
// @Param uuid path string true "User Uuid"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/users/{uuid}/block [put]
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

// UnblockUser godoc
// @Summary Unblock a user
// @Description Unblock a user by their Uuid. Only admin can perform this action.
// @Tags Auth
// @Accept json
// @Produce json
// @Param uuid path string true "User Uuid"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/users/{uuid}/unblock [put]
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