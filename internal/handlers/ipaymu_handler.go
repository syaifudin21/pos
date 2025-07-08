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
	"github.com/msyaifudin/pos/internal/validators"
)

var validate = validator.New()

type IpaymuHandler struct {
	Service *services.IpaymuService
}

func NewIpaymuHandler(service *services.IpaymuService) *IpaymuHandler {
	return &IpaymuHandler{Service: service}
}

func (h *IpaymuHandler) CreateDirectPayment(c echo.Context) error {
	var req dtos.CreateDirectPaymentRequest
	if err := c.Bind(&req); err != nil {
		// Check if it's a binding error (e.g., JSON parsing, type mismatch)
		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusBadRequest {
			return JSONError(c, http.StatusBadRequest, "Invalid JSON format or data type mismatch.")
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
		// Check if it's a binding error (e.g., JSON parsing, type mismatch)
		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusBadRequest {
			return JSONError(c, http.StatusBadRequest, "Invalid JSON format or data type mismatch.")
		}
		return JSONError(c, http.StatusBadRequest, "invalid_input")
	}

	// Validate the request struct
	if err := validate.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make([]string, 0, len(ve))
			for _, fe := range ve {
				errMsgs = append(errMsgs, fmt.Sprintf("ssss '%s' failed on the '%s' tag", fe.Field(), fe.Tag()))
			}
			return JSONError(c, http.StatusBadRequest, strings.Join(errMsgs, ", "))
		}
		return JSONError(c, http.StatusBadRequest, err.Error())
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
		// Check if it's a binding error (e.g., JSON parsing, type mismatch)
		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusBadRequest {
			return JSONError(c, http.StatusBadRequest, "Invalid JSON format or data type mismatch.")
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
