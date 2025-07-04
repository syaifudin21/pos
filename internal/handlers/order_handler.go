package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/pkg/utils"
)

type OrderHandler struct {
	OrderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: orderService}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with specified products and quantities.
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order details"
// @Success 201 {object} SuccessResponse{data=models.Order}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	req := new(models.CreateOrderRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	claims := c.Get("claims").(*utils.Claims)
	userUuid := claims.ID // Assuming user's uuid is stored in claims.ID

	// For simplicity, assuming outlet ID is passed in the request or derived from user's outlet
	// For now, let's use the outlet ID from the request
	outletUuid := req.OutletUuid

	order, err := h.OrderService.CreateOrder(outletUuid, userUuid, req.Items)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{Message: "Order created successfully", Data: order})
}

// GetOrderByUuid godoc
// @Summary Get order by Uuid
// @Description Get a single order by its Uuid.
// @Tags Orders
// @Accept json
// @Produce json
// @Param uuid path string true "Order Uuid"
// @Success 200 {object} SuccessResponse{data=models.Order}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{uuid} [get]
func (h *OrderHandler) GetOrderByUuid(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid UUID format"})
	}

	order, err := h.OrderService.GetOrderByUuid(uuid)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Order retrieved successfully", Data: order})
}

// GetOrdersByOutlet godoc
// @Summary Get all orders for an outlet
// @Description Get a list of all orders for a given outlet.
// @Tags Orders
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet Uuid"
// @Success 200 {object} SuccessResponse{data=[]models.Order}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/orders [get]
func (h *OrderHandler) GetOrdersByOutlet(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Outlet Uuid format"})
	}

	orders, err := h.OrderService.GetOrdersByOutlet(outletUuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Orders retrieved successfully", Data: orders})
}
