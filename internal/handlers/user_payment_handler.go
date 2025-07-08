package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

var validate = validator.New()

type UserPaymentHandler struct {
	UserPaymentService *services.UserPaymentService
	UserContextService *services.UserContextService
}

func NewUserPaymentHandler(userPaymentService *services.UserPaymentService, userContextService *services.UserContextService) *UserPaymentHandler {
	return &UserPaymentHandler{
		UserPaymentService: userPaymentService,
		UserContextService: userContextService,
	}
}

func (h *UserPaymentHandler) ActivateUserPayment(c echo.Context) error {
	var req dtos.ActivateUserPaymentRequest
	if err := c.Bind(&req); err != nil {
		// Check if it's a binding error (e.g., JSON parsing, type mismatch)
		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusBadRequest {
			return JSONError(c, http.StatusBadRequest, "Invalid JSON format or data type mismatch. Please ensure payment_method_id is an integer.")
		}
		return JSONError(c, http.StatusBadRequest, "invalid_input")
	}

	// Validate the request struct
	if err := validate.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make([]string, 0, len(ve))
			for _, fe := range ve {
				errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' failed on the '%s' tag", fe.Field(), fe.Tag()))
			}
			return JSONError(c, http.StatusBadRequest, strings.Join(errMsgs, ", "))
		}
		return JSONError(c, http.StatusBadRequest, err.Error())
	}

	// Get UserID from context (set by SelfAuthorize middleware)
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_user_id")
	}

	err = h.UserPaymentService.ActivateUserPayment(userID, req.PaymentMethodID)
	if err != nil {
		if errors.Is(err, services.ErrIpaymuRegistrationRequired) {
			return c.JSON(http.StatusPreconditionRequired, map[string]string{
				"message": "iPaymu registration required",
				"route":   "/ipaymu/register",
			})
		} else if errors.Is(err, services.ErrTsmRegistrationRequired) {
			return c.JSON(http.StatusPreconditionRequired, map[string]string{
				"message": "TSM registration required",
				"route":   "/tsm/register",
			})
		}
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "payment_method_activated_successfully", nil)
}

func (h *UserPaymentHandler) DeactivateUserPayment(c echo.Context) error {
	var req dtos.DeactivateUserPaymentRequest
	if err := c.Bind(&req); err != nil {
		// Check if it's a binding error (e.g., JSON parsing, type mismatch)
		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusBadRequest {
			return JSONError(c, http.StatusBadRequest, "Invalid JSON format or data type mismatch. Please ensure payment_method_id is an integer.")
		}
		return JSONError(c, http.StatusBadRequest, "invalid_input")
	}

	// Validate the request struct
	if err := validate.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make([]string, 0, len(ve))
			for _, fe := range ve {
				errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' failed on the '%s' tag", fe.Field(), fe.Tag()))
			}
			return JSONError(c, http.StatusBadRequest, strings.Join(errMsgs, ", "))
		}
		return JSONError(c, http.StatusBadRequest, err.Error())
	}

	// Get UserID from context (set by SelfAuthorize middleware)
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_user_id")
	}

	if err := h.UserPaymentService.DeactivateUserPayment(userID, req.PaymentMethodID); err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "payment_method_deactivated_successfully", nil)
}

func (h *UserPaymentHandler) ListUserPayments(c echo.Context) error {
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_user_id")
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_owner_id")
	}

	userPayments, err := h.UserPaymentService.ListUserPaymentsByOwner(ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	var responses []dtos.UserPaymentResponse
	for _, up := range userPayments {
		responses = append(responses, dtos.UserPaymentResponse{
			PaymentMethodID: up.PaymentMethodID,
			PaymentName:     up.PaymentMethod.Name,
			PaymentMethod:   up.PaymentMethod.PaymentMethod,
			PaymentChannel:  up.PaymentMethod.PaymentChannel,
			IsActive:        up.IsActive,
		})
	}

	return JSONSuccess(c, http.StatusOK, "user_payments_listed_successfully", responses)
}