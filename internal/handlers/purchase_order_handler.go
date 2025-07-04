package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/services"
)

type PurchaseOrderHandler struct {
	PurchaseOrderService *services.PurchaseOrderService
}

func NewPurchaseOrderHandler(poService *services.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{PurchaseOrderService: poService}
}

// CreatePurchaseOrder godoc
// @Summary Create a new purchase order
// @Description Create a new purchase order with specified supplier, outlet, and products.
// @Tags Purchase Orders
// @Accept json
// @Produce json
// @Param purchase_order body models.CreatePurchaseOrderRequest true "Purchase order details"
// @Success 201 {object} SuccessResponse{data=models.PurchaseOrder}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /purchase-orders [post]
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c echo.Context) error {
	req := new(models.CreatePurchaseOrderRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	po, err := h.PurchaseOrderService.CreatePurchaseOrder(req.SupplierUuid, req.OutletUuid, req.Items)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{Message: "Purchase order created successfully", Data: po})
}

// GetPurchaseOrderByUuid godoc
// @Summary Get purchase order by Uuid
// @Description Get a single purchase order by its Uuid.
// @Tags Purchase Orders
// @Accept json
// @Produce json
// @Param uuid path string true "Purchase Order Uuid"
// @Success 200 {object} SuccessResponse{data=models.PurchaseOrder}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /purchase-orders/{uuid} [get]
func (h *PurchaseOrderHandler) GetPurchaseOrderByUuid(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Uuid format"})
	}

	po, err := h.PurchaseOrderService.GetPurchaseOrderByUuid(uuid)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Purchase order retrieved successfully", Data: po})
}

// GetPurchaseOrdersByOutlet godoc
// @Summary Get all purchase orders for an outlet
// @Description Get a list of all purchase orders for a given outlet.
// @Tags Purchase Orders
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet Uuid"
// @Success 200 {object} SuccessResponse{data=[]models.PurchaseOrder}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/purchase-orders [get]
func (h *PurchaseOrderHandler) GetPurchaseOrdersByOutlet(c echo.Context) error {
	OutletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Outlet Uuid format"})
	}

	pos, err := h.PurchaseOrderService.GetPurchaseOrdersByOutlet(OutletUuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Purchase orders retrieved successfully", Data: pos})
}

// ReceivePurchaseOrder godoc
// @Summary Receive a purchase order
// @Description Mark a purchase order as completed and update stock quantities.
// @Tags Purchase Orders
// @Accept json
// @Produce json
// @Param uuid path string true "Purchase Order Uuid"
// @Success 200 {object} SuccessResponse{data=models.PurchaseOrder}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /purchase-orders/{uuid}/receive [put]
func (h *PurchaseOrderHandler) ReceivePurchaseOrder(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Uuid format"})
	}

	po, err := h.PurchaseOrderService.ReceivePurchaseOrder(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Purchase order received successfully", Data: po})
}
