package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
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
// @Param outlet_uuid path string true "Outlet Uuid"
// @Param product_uuid path string true "Product Uuid"
// @Success 200 {object} SuccessResponse{data=dtos.StockResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/stocks/{product_uuid} [get]
func (h *StockHandler) GetStockByOutletAndProduct(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_product_uuid_format")
	}

	stock, err := h.StockService.GetStockByOutletAndProduct(outletUuid, productUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "stock_retrieved_successfully", stock)
}

// GetOutletStocks godoc
// @Summary Get all stocks for an outlet
// @Description Get a list of all stock quantities for products in a given outlet.
// @Tags Stocks
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet Uuid"
// @Success 200 {object} SuccessResponse{data=[]dtos.StockResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/stocks [get]
func (h *StockHandler) GetOutletStocks(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	stocks, err := h.StockService.GetOutletStocks(outletUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "outlet_stocks_retrieved_successfully", stocks)
}

// UpdateStock godoc
// @Summary Update stock quantity
// @Description Update the stock quantity for a specific product in an outlet. This is for direct stock adjustments.
// @Tags Stocks
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet Uuid"
// @Param product_uuid path string true "Product Uuid"
// @Success 200 {object} SuccessResponse{data=dtos.StockResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/stocks/{product_uuid} [put]
func (h *StockHandler) UpdateStock(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_product_uuid_format")
	}

	req := new(dtos.UpdateStockRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	stock, err := h.StockService.UpdateStock(outletUuid, productUuid, req.Quantity)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "stock_updated_successfully", stock)
}

// UpdateGlobalStock godoc
// @Summary Update stock quantity globally
// @Description Update the stock quantity for a specific product in an outlet by providing IDs in the request body.
// @Tags Stocks
// @Accept json
// @Produce json
// @Param stock body dtos.GlobalStockUpdateRequest true "Stock update details"
// @Success 200 {object} SuccessResponse{data=dtos.StockResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /stocks [put]
func (h *StockHandler) UpdateGlobalStock(c echo.Context) error {
	req := new(dtos.GlobalStockUpdateRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	stock, err := h.StockService.UpdateGlobalStock(req.OutletUuid, req.Productuuid, req.Quantity)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "stock_updated_successfully", stock)
}