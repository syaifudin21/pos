package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type PurchaseOrderHandler struct {
	PurchaseOrderService *services.PurchaseOrderService
	UserContextService   *services.UserContextService
}

func NewPurchaseOrderHandler(poService *services.PurchaseOrderService, userContextService *services.UserContextService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{PurchaseOrderService: poService, UserContextService: userContextService}
}
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.CreatePurchaseOrderRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	po, err := h.PurchaseOrderService.CreatePurchaseOrder(*req, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "purchase_order_created_successfully", po)
}

func (h *PurchaseOrderHandler) GetPurchaseOrderByUuid(c echo.Context) error {
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

	po, err := h.PurchaseOrderService.GetPurchaseOrderByUuid(parsedUuid, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "purchase_order_retrieved_successfully", po)
}

func (h *PurchaseOrderHandler) GetPurchaseOrdersByOutlet(c echo.Context) error {
	OutletUuidParam := c.Param("outlet_uuid")
	OutletUuid, err := uuid.Parse(OutletUuidParam)
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

	pos, err := h.PurchaseOrderService.GetPurchaseOrdersByOutlet(OutletUuid, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "purchase_orders_retrieved_successfully", pos)
}

func (h *PurchaseOrderHandler) ReceivePurchaseOrder(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	po, err := h.PurchaseOrderService.ReceivePurchaseOrder(parsedUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "purchase_order_received_successfully", po)
}
