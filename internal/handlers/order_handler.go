package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type OrderHandler struct {
	OrderService       *services.OrderService
	UserContextService *services.UserContextService
}

func NewOrderHandler(orderService *services.OrderService, userContextService *services.UserContextService) *OrderHandler {
	return &OrderHandler{OrderService: orderService, UserContextService: userContextService}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.CreateOrderRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	order, err := h.OrderService.CreateOrder(*req, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "order_created_successfully", order)
}

func (h *OrderHandler) GetOrderByUuid(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	order, err := h.OrderService.GetOrderByUuid(parsedUuid, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "order_retrieved_successfully", order)
}
func (h *OrderHandler) GetOrdersByOutlet(c echo.Context) error {
	outletUuidParam := c.Param("outlet_uuid")
	outletUuid, err := uuid.Parse(outletUuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	orders, err := h.OrderService.GetOrdersByOutlet(outletUuid, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "orders_retrieved_successfully", orders)
}
