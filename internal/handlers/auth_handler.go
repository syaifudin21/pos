package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/services"
)

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
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} SuccessResponse{data=UserResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	if req.Username == "" || req.Password == "" || req.Role == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Username, password, and role are required"})
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Password, req.Role, req.OutletID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{Message: "User registered successfully", Data: UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role, OutletID: user.OutletID}})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate user and return a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} SuccessResponse{data=LoginResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Username and password are required"})
	}

	token, user, err := h.AuthService.LoginUser(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, SuccessResponse{Message: "Login successful", Data: LoginResponse{Token: token, User: UserResponse{ID: user.ID, Uuid: user.Uuid, Username: user.Username, Role: user.Role, OutletID: user.OutletID}}})
}

// Request and Response Structs for Swagger
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	OutletID *uint  `json:"outlet_id,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       uint      `json:"id"`
	Uuid     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	OutletID *uint     `json:"outlet_id,omitempty"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// BlockUser godoc
// @Summary Block a user
// @Description Block a user by their External ID. Only admin can perform this action.
// @Tags Auth
// @Accept json
// @Produce json
// @Param uuid path string true "User External ID (UUID)"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/users/{uuid}/block [put]
func (h *AuthHandler) BlockUser(c echo.Context) error {
	userExternalID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid User External ID format"})
	}

	if err := h.AuthService.BlockUser(userExternalID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, SuccessResponse{Message: "User blocked successfully"})
}

// UnblockUser godoc
// @Summary Unblock a user
// @Description Unblock a user by their External ID. Only admin can perform this action.
// @Tags Auth
// @Accept json
// @Produce json
// @Param uuid path string true "User External ID (UUID)"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/users/{uuid}/unblock [put]
func (h *AuthHandler) UnblockUser(c echo.Context) error {
	userExternalID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid User External ID format"})
	}

	if err := h.AuthService.UnblockUser(userExternalID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, SuccessResponse{Message: "User unblocked successfully"})
}
