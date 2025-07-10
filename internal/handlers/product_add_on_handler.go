package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type ProductAddOnHandler struct {
	ProductAddOnService *services.ProductAddOnService
	UserContextService  *services.UserContextService
}

func NewProductAddOnHandler(productAddOnService *services.ProductAddOnService, userContextService *services.UserContextService) *ProductAddOnHandler {
	return &ProductAddOnHandler{ProductAddOnService: productAddOnService, UserContextService: userContextService}
}

func (h *ProductAddOnHandler) CreateProductAddOn(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.ProductAddOnRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	resp, err := h.ProductAddOnService.CreateProductAddOn(req, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusCreated, "product_add_on_created_successfully", resp)
}

func (h *ProductAddOnHandler) GetProductAddOnsByProductID(c echo.Context) error {
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_product_uuid_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	resp, err := h.ProductAddOnService.GetProductAddOnsByProductID(productUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "product_add_ons_retrieved_successfully", resp)
}

func (h *ProductAddOnHandler) DeleteProductAddOn(c echo.Context) error {
	productAddOnUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	err = h.ProductAddOnService.DeleteProductAddOn(productAddOnUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "product_add_on_deleted_successfully", nil)
}
