package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type OrderPaymentHandler struct {
	OrderPaymentService *services.OrderPaymentService
	UserContextService  *services.UserContextService
}

func NewOrderPaymentHandler(orderPaymentService *services.OrderPaymentService, userContextService *services.UserContextService) *OrderPaymentHandler {
	return &OrderPaymentHandler{OrderPaymentService: orderPaymentService, UserContextService: userContextService}
}

func (h *OrderPaymentHandler) CreateOrderPayment(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.CreateOrderPaymentRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	orderPayment, err := h.OrderPaymentService.CreateOrderPayment(*req, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "order_payment_created_successfully", orderPayment)
}
