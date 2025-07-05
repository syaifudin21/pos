package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/validators"
	"github.com/msyaifudin/pos/pkg/utils"
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

	lang := c.Get("lang").(string)
	if messages := validators.ValidateCreatePurchaseOrder(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.PurchaseOrderService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	po, err := h.PurchaseOrderService.CreatePurchaseOrder(req.SupplierUuid, req.OutletUuid, req.Items, ownerID)
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.PurchaseOrderService.GetOwnerID(userID)
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.PurchaseOrderService.GetOwnerID(userID)
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	po, err := h.PurchaseOrderService.ReceivePurchaseOrder(parsedUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "purchase_order_received_successfully", po)
}
