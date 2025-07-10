package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type StockHandler struct {
	StockService       *services.StockService
	UserContextService *services.UserContextService
}

func NewStockHandler(stockService *services.StockService, userContextService *services.UserContextService) *StockHandler {
	return &StockHandler{StockService: stockService, UserContextService: userContextService}
}

func (h *StockHandler) GetOutletStocks(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	productType := c.QueryParam("product_type")
	isForSaleStr := c.QueryParam("is_for_sale")
	isForSale := false
	if isForSaleStr == "true" {
		isForSale = true
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	stocks, err := h.StockService.GetOutletStocks(outletUuid, userID, productType, isForSale)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "outlet_stocks_retrieved_successfully", stocks)
}

func (h *StockHandler) UpdateStock(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	req, ok := c.Get("validated_data").(*dtos.UpdateStockRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	stock, err := h.StockService.UpdateStock(*req, outletUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "stock_updated_successfully", stock)
}

