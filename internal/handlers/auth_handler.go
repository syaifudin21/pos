package handlers

import (
	"log"
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

	user, err := h.AuthService.RegisterUser(req.Name, req.Email, req.Password, req.Role, req.OutletID, &creatorID, nil, false)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "user_registered_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
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
	user, err := h.AuthService.RegisterUser(req.Name, req.Email, req.Password, "admin", nil, nil, &req.PhoneNumber, false)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "admin_registered_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
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

	token, user, err := h.AuthService.LoginUser(req.Email, req.Password)

	if err != nil {
		statusCode := MapErrorToStatusCode(err)

		if err.Error() == "user not verified" {
			return JSONError(c, statusCode, err.Error())
		}
		return JSONError(c, statusCode, err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "login_successful", dtos.LoginResponse{Token: token, User: dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role}})
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
		userResponses = append(userResponses, dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
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

	return JSONSuccess(c, http.StatusOK, "otp_verified_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	claims := c.Get("user").(*jwt.Token).Claims.(*utils.Claims)
	userID := claims.ID

	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "profile_retrieved_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
}

func (h *AuthHandler) UpdatePassword(c echo.Context) error {
	claims := c.Get("user").(*jwt.Token).Claims.(*utils.Claims)
	userID := claims.ID

	req := new(dtos.UpdatePasswordRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("UpdatePassword Handler: Invalid request payload: %v", err)
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	log.Printf("UpdatePassword Handler: UserID: %d, OldPassword: %s, NewPassword: %s", userID, req.OldPassword, req.NewPassword)

	lang := c.Get("lang").(string)
	if messages := validators.ValidateUpdatePasswordRequest(req, lang); messages != nil {
		log.Printf("UpdatePassword Handler: Validation failed: %v", messages)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	if err := h.AuthService.UpdatePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		log.Printf("UpdatePassword Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("UpdatePassword Handler: Password updated successfully")
	return JSONSuccess(c, http.StatusOK, "password_updated_successfully", nil)
}

func (h *AuthHandler) SendOTPForEmailUpdate(c echo.Context) error {
	claims := c.Get("user").(*jwt.Token).Claims.(*utils.Claims)
	userID := claims.ID

	req := new(dtos.SendOTPRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateSendOTPRequest(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	if err := h.AuthService.SendOTPForEmailUpdate(userID, req.Email); err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "otp_sent_for_email_update", nil)
}

// Fungsi ini seharusnya bisa mengambil request body jika request yang dikirimkan bertipe JSON dan field-field pada body sesuai dengan struct dtos.UpdateEmailRequest.
// Namun, jika req selalu kosong setelah c.Bind(req), kemungkinan penyebabnya adalah:
// 1. Header Content-Type pada request tidak di-set ke "application/json".
// 2. Field pada JSON body tidak sesuai dengan tag json di struct UpdateEmailRequest (harus "new_email" dan "otp").
// 3. Body request kosong atau tidak valid JSON.
// 4. Ada middleware yang memodifikasi body sebelum sampai ke handler.

// Berikut contoh rewrite dengan log tambahan untuk membantu debug:
func (h *AuthHandler) UpdateEmail(c echo.Context) error {
	claims := c.Get("user").(*jwt.Token).Claims.(*utils.Claims)
	userID := claims.ID

	req := new(dtos.UpdateEmailRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("UpdateEmail Handler: Invalid request payload: %v", err)
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	log.Printf("UpdateEmail Handler: UserID: %d, NewEmail: %s, OTP: %s", userID, req.NewEmail, req.OTP)

	lang := c.Get("lang").(string)
	if messages := validators.ValidateUpdateEmailRequest(req, lang); messages != nil {
		log.Printf("UpdateEmail Handler: Validation failed: %v", messages)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	if err := h.AuthService.UpdateEmail(userID, req.NewEmail, req.OTP); err != nil {
		log.Printf("UpdateEmail Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("UpdateEmail Handler: Email updated successfully")
	return JSONSuccess(c, http.StatusOK, "email_updated_successfully", nil)
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	req := new(dtos.ForgotPasswordRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("ForgotPassword Handler: Invalid request payload: %v", err)
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	log.Printf("ForgotPassword Handler: Email: %s", req.Email)

	lang := c.Get("lang").(string)
	if messages := validators.ValidateForgotPasswordRequest(req, lang); messages != nil {
		log.Printf("ForgotPassword Handler: Validation failed: %v", messages)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	if err := h.AuthService.SendOTPForPasswordReset(req.Email); err != nil {
		log.Printf("ForgotPassword Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("ForgotPassword Handler: OTP sent for password reset")
	return JSONSuccess(c, http.StatusOK, "otp_sent_for_password_reset", nil)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	req := new(dtos.ResetPasswordRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("ResetPassword Handler: Invalid request payload: %v", err)
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	log.Printf("ResetPassword Handler: Email: %s, OTP: %s, NewPassword: %s", req.Email, req.OTP, req.NewPassword)

	lang := c.Get("lang").(string)
	if messages := validators.ValidateResetPasswordRequest(req, lang); messages != nil {
		log.Printf("ResetPassword Handler: Validation failed: %v", messages)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	if err := h.AuthService.ResetPassword(req.Email, req.OTP, req.NewPassword); err != nil {
		log.Printf("ResetPassword Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("ResetPassword Handler: Password reset successfully")
	return JSONSuccess(c, http.StatusOK, "password_reset_successful", nil)
}
