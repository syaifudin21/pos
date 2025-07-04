package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/services"
)

type StockHandler struct {
	StockService *services.StockService
}

func NewStockHandler(stockService *services.StockService) *StockHandler {
	return &StockHandler{StockService: stockService}
}

// GetStockByOutletAndProduct godoc
// @Summary Get stock for a product in an outlet
// @Description Get the stock quantity for a specific product in a given outlet.
// @Tags Stocks
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet External ID (UUID)"
// @Param product_uuid path string true "Product External ID (UUID)"
// @Success 200 {object} SuccessResponse{data=models.Stock}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/stocks/{product_uuid} [get]
func (h *StockHandler) GetStockByOutletAndProduct(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Outlet External ID format"})
	}
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Product External ID format"})
	}

	stock, err := h.StockService.GetStockByOutletAndProduct(outletUuid, productUuid)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Stock retrieved successfully", Data: stock})
}

// GetOutletStocks godoc
// @Summary Get all stocks for an outlet
// @Description Get a list of all stock quantities for products in a given outlet.
// @Tags Stocks
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet External ID (UUID)"
// @Success 200 {object} SuccessResponse{data=[]models.Stock}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/stocks [get]
func (h *StockHandler) GetOutletStocks(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Outlet External ID format"})
	}

	stocks, err := h.StockService.GetOutletStocks(outletUuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Outlet stocks retrieved successfully", Data: stocks})
}

// UpdateStock godoc
// @Summary Update stock quantity
// @Description Update the stock quantity for a specific product in an outlet. This is for direct stock adjustments.
// @Tags Stocks
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet External ID (UUID)"
// @Param product_uuid path string true "Product External ID (UUID)"
// @Param stock body UpdateStockRequest true "Stock update details"
// @Success 200 {object} SuccessResponse{data=models.Stock}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/stocks/{product_uuid} [put]
func (h *StockHandler) UpdateStock(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Outlet External ID format"})
	}
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Product External ID format"})
	}

	req := new(UpdateStockRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	stock, err := h.StockService.UpdateStock(outletUuid, productUuid, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Stock updated successfully", Data: stock})
}

type UpdateStockRequest struct {
	Quantity float64 `json:"quantity"`
}
