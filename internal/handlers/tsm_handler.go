package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type TsmHandler struct {
	TsmService         *services.TsmService
	UserContextService *services.UserContextService
	UserPaymentService *services.UserPaymentService
}

func NewTsmHandler(tsmService *services.TsmService, userContextService *services.UserContextService, userPaymentService *services.UserPaymentService) *TsmHandler {
	return &TsmHandler{
		TsmService:         tsmService,
		UserContextService: userContextService,
		UserPaymentService: userPaymentService,
	}
}

func (h *TsmHandler) RegisterTsm(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.TsmRegisterRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_user_id")
	}

	hasIpaymu, err := h.UserPaymentService.HasIpaymuConnection(userID)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "failed_to_check_ipaymu_connection")
	}
	if !hasIpaymu {
		return JSONError(c, http.StatusForbidden, "ipaymu_registration_required_to_access_tsm_register")
	}

	if err := h.TsmService.RegisterTsm(userID, *req); err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "tsm_registered_successfully", nil)
}
