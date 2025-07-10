package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type IpaymuHandler struct {
	Service            *services.IpaymuService
	UserContextService *services.UserContextService
}

func NewIpaymuHandler(service *services.IpaymuService, userContextService *services.UserContextService) *IpaymuHandler {
	return &IpaymuHandler{Service: service, UserContextService: userContextService}
}

func (h *IpaymuHandler) CreateDirectPayment(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.CreateDirectPaymentRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_user_id")
	}

	res, err := h.Service.CreateDirectPayment(
		userID,
		req.ServiceName, req.ServiceRefID, req.Product, req.Qty, req.Price, req.Name, req.Email, req.Phone, req.Method, req.Channel, req.Account,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// IpaymuNotify handles notify/callback from Ipaymu
func (h *IpaymuHandler) IpaymuNotify(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.IpaymuNotifyRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	// fmt.Println("Debug Signature:", req) // Debugging output

	err := h.Service.NotifyDirectPayment(req.TrxID, req.Status, req.SettlementStatus)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Success"})
}

// Handler untuk register user ke Ipaymu
func (h *IpaymuHandler) RegisterIpaymu(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.RegisterIpaymuRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	optional := make(map[string]interface{})
	if req.IdentityNo != nil {
		optional["identityNo"] = *req.IdentityNo
	}
	if req.BusinessName != nil {
		optional["businessName"] = *req.BusinessName
	}
	if req.Birthday != nil {
		optional["birthday"] = *req.Birthday
	}
	if req.Birthplace != nil {
		optional["birthplace"] = *req.Birthplace
	}
	if req.Gender != nil {
		optional["gender"] = *req.Gender
	}
	if req.Address != nil {
		optional["address"] = *req.Address
	}
	if req.IdentityPhoto != nil {
		optional["identityPhoto"] = req.IdentityPhoto
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_user_id")
	}

	res, err := h.Service.Register(
		userID,
		req.Name,
		req.Phone,
		req.Password,
		req.Email,
		optional,
	)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "ipaymu_registration_successful", res)
}
