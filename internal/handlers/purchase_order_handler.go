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
}

func NewPurchaseOrderHandler(poService *services.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{PurchaseOrderService: poService}
}
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c echo.Context) error {
	req := new(dtos.CreatePurchaseOrderRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	po, err := h.PurchaseOrderService.CreatePurchaseOrder(req.SupplierUuid, req.OutletUuid, req.Items)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "purchase_order_created_successfully", po)
}

func (h *PurchaseOrderHandler) GetPurchaseOrderByUuid(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	po, err := h.PurchaseOrderService.GetPurchaseOrderByUuid(uuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "purchase_order_retrieved_successfully", po)
}

func (h *PurchaseOrderHandler) GetPurchaseOrdersByOutlet(c echo.Context) error {
	OutletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	pos, err := h.PurchaseOrderService.GetPurchaseOrdersByOutlet(OutletUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "purchase_orders_retrieved_successfully", pos)
}

func (h *PurchaseOrderHandler) ReceivePurchaseOrder(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	po, err := h.PurchaseOrderService.ReceivePurchaseOrder(uuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "purchase_order_received_successfully", po)
}
