package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/validators"
)

type IpaymuHandler struct {
	Service *services.IpaymuService
}

func NewIpaymuHandler(service *services.IpaymuService) *IpaymuHandler {
	return &IpaymuHandler{Service: service}
}

func (h *IpaymuHandler) CreateDirectPayment(c echo.Context) error {
	var req dtos.CreateDirectPaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request", "details": err.Error()})
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateCreateDirectPayment(&req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	res, err := h.Service.CreateDirectPayment(
		req.ServiceName, req.ServiceRefID, req.Product, req.Qty, req.Price, req.Name, req.Email, req.Phone, req.Method, req.Channel, req.Account,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// IpaymuNotify handles notify/callback from Ipaymu
func (h *IpaymuHandler) IpaymuNotify(c echo.Context) error {
	var req dtos.IpaymuNotifyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request", "details": err.Error()})
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateIpaymuNotify(&req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}
	fmt.Println("Debug Signature:", req) // Debugging output

	err := h.Service.NotifyDirectPayment(req.TrxID, req.Status, req.SettlementStatus)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Success"})
}

// Handler untuk register user ke Ipaymu
func (h *IpaymuHandler) RegisterIpaymu(c echo.Context) error {
	var req dtos.RegisterIpaymuRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request", "details": err.Error()})
	}

	// Validasi opsional, bisa tambahkan validator custom jika perlu
	if req.Name == "" || req.Phone == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "name, phone, dan password wajib diisi"})
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

	res, err := h.Service.Register(
		req.Name,
		req.Phone,
		req.Password,
		req.Email,
		optional,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
