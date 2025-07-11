package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type AuthHandler struct {
	AuthService        *services.AuthService
	UserContextService *services.UserContextService
}

func NewAuthHandler(authService *services.AuthService, userContextService *services.UserContextService) *AuthHandler {
	return &AuthHandler{AuthService: authService, UserContextService: userContextService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.RegisterRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	// Get the ID of the currently logged-in admin from the JWT claims
	creatorID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	user, err := h.AuthService.RegisterUser(req.Name, req.Email, req.Password, req.Role, req.OutletID, &creatorID, nil, false)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "user_registered_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
}

func (h *AuthHandler) RegisterOwner(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.RegisterOwnerRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	// This endpoint is for public owner registration, so no creatorID is needed.
	// The role is hardcoded to "owner".
	user, err := h.AuthService.RegisterUser(req.Name, req.Email, req.Password, "owner", nil, nil, &req.PhoneNumber, false)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "owner_registered_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.LoginRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
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
	adminID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

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
	req, ok := c.Get("validated_data").(*dtos.UpdateUserRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	user, err := h.AuthService.GetUserByuuid(userUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	if err := h.AuthService.UpdateUser(user.ID, req); err != nil {
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
	req, ok := c.Get("validated_data").(*dtos.VerifyOTPRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	user, err := h.AuthService.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "otp_verified_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "profile_retrieved_successfully", dtos.UserResponse{ID: user.ID, Uuid: user.Uuid, Name: user.Name, Email: user.Email, Role: user.Role})
}

func (h *AuthHandler) UpdatePassword(c echo.Context) error {
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	req, ok := c.Get("validated_data").(*dtos.UpdatePasswordRequest)
	if !ok {
		log.Printf("UpdatePassword Handler: Failed to get validated request from context.")
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	log.Printf("UpdatePassword Handler: UserID: %d, OldPassword: %s, NewPassword: %s", userID, req.OldPassword, req.NewPassword)

	if err := h.AuthService.UpdatePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		log.Printf("UpdatePassword Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("UpdatePassword Handler: Password updated successfully")
	return JSONSuccess(c, http.StatusOK, "password_updated_successfully", nil)
}

func (h *AuthHandler) SendOTPForEmailUpdate(c echo.Context) error {
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	req, ok := c.Get("validated_data").(*dtos.SendOTPRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
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
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	req, ok := c.Get("validated_data").(*dtos.UpdateEmailRequest)
	if !ok {
		log.Printf("UpdateEmail Handler: Failed to get validated request from context.")
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	log.Printf("UpdateEmail Handler: UserID: %d, NewEmail: %s, OTP: %s", userID, req.NewEmail, req.OTP)

	if err := h.AuthService.UpdateEmail(userID, req.NewEmail, req.OTP); err != nil {
		log.Printf("UpdateEmail Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("UpdateEmail Handler: Email updated successfully")
	return JSONSuccess(c, http.StatusOK, "email_updated_successfully", nil)
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.ForgotPasswordRequest)
	if !ok {
		log.Printf("ForgotPassword Handler: Failed to get validated request from context.")
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	log.Printf("ForgotPassword Handler: Email: %s", req.Email)

	if err := h.AuthService.SendOTPForPasswordReset(req.Email); err != nil {
		log.Printf("ForgotPassword Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("ForgotPassword Handler: OTP sent for password reset")
	return JSONSuccess(c, http.StatusOK, "otp_sent_for_password_reset", nil)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.ResetPasswordRequest)
	if !ok {
		log.Printf("ResetPassword Handler: Failed to get validated request from context.")
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	log.Printf("ResetPassword Handler: Email: %s, OTP: %s, NewPassword: %s", req.Email, req.OTP, req.NewPassword)

	if err := h.AuthService.ResetPassword(req.Email, req.OTP, req.NewPassword); err != nil {
		log.Printf("ResetPassword Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("ResetPassword Handler: Password reset successful")
	return JSONSuccess(c, http.StatusOK, "password_reset_successful", nil)
}

func (h *AuthHandler) ResendVerificationEmail(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.ResendEmailRequest)
	if !ok {
		log.Printf("ResendVerificationEmail Handler: Failed to get validated request from context.")
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	log.Printf("ResendVerificationEmail Handler: Email: %s", req.Email)

	if err := h.AuthService.ResendVerificationEmail(req.Email); err != nil {
		log.Printf("ResendVerificationEmail Handler: Service error: %v", err)
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	log.Println("ResendVerificationEmail Handler: Verification email resent successfully")
	return JSONSuccess(c, http.StatusOK, "verification_email_resent_successfully", nil)
}
