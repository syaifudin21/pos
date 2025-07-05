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

type StockHandler struct {
	StockService *services.StockService
}

func NewStockHandler(stockService *services.StockService) *StockHandler {
	return &StockHandler{StockService: stockService}
}

func (h *StockHandler) GetStockByOutletAndProduct(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_product_uuid_format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.StockService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	stock, err := h.StockService.GetStockByOutletAndProduct(outletUuid, productUuid, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "stock_retrieved_successfully", stock)
}

func (h *StockHandler) GetOutletStocks(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.StockService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	stocks, err := h.StockService.GetOutletStocks(outletUuid, ownerID)
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
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_product_uuid_format")
	}

	req := new(dtos.UpdateStockRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateUpdateStock(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.StockService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	stock, err := h.StockService.UpdateStock(outletUuid, productUuid, req.Quantity, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "stock_updated_successfully", stock)
}

func (h *StockHandler) UpdateGlobalStock(c echo.Context) error {
	req := new(dtos.GlobalStockUpdateRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateGlobalStockUpdate(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.StockService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	stock, err := h.StockService.UpdateGlobalStock(req.OutletUuid, req.Productuuid, req.Quantity, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "stock_updated_successfully", stock)
}
