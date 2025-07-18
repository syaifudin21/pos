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
		return JSONError(c, http.StatusBadRequest, "invalid_request_body")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, "unauthorized")
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, http.StatusForbidden, "forbidden")
	}

	if err := h.TsmService.RegisterTsm(ownerID, *req); err != nil {
		return JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "tsm_registered_successfully", nil)
}

func (h *TsmHandler) GenerateApplink(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.TsmGenerateApplinkRequest)
	if !ok {
		return JSONError(c, http.StatusBadRequest, "invalid_request_body")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, "unauthorized")
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, http.StatusForbidden, "forbidden")
	}

	resp, err := h.TsmService.GenerateAPPLink(ownerID, *req)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "applink_generated_successfully", resp)
}

func (h *TsmHandler) Callback(c echo.Context) error {
	var req dtos.TsmCallbackRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_body")
	}

	if err := h.TsmService.HandleCallback(req); err != nil {
		return JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "callback_processed_successfully", nil)
}
