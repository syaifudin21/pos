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

type SupplierHandler struct {
	SupplierService *services.SupplierService
}

func NewSupplierHandler(supplierService *services.SupplierService) *SupplierHandler {
	return &SupplierHandler{SupplierService: supplierService}
}

func (h *SupplierHandler) GetAllSuppliers(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.SupplierService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	suppliers, err := h.SupplierService.GetAllSuppliers(ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	var supplierResponses []dtos.SupplierResponse
	for _, supplier := range suppliers {
		supplierResponses = append(supplierResponses, dtos.SupplierResponse{
			ID:      supplier.ID,
			Uuid:    supplier.Uuid,
			Name:    supplier.Name,
			Contact: supplier.Contact,
			Address: supplier.Address,
		})
	}
	return JSONSuccess(c, http.StatusOK, "Suppliers retrieved successfully", supplierResponses)
}
func (h *SupplierHandler) GetSupplierByuuid(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid Uuid format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.SupplierService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	supplier, err := h.SupplierService.GetSupplierByuuid(parsedUuid, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "Supplier retrieved successfully", supplier)
}

func (h *SupplierHandler) CreateSupplier(c echo.Context) error {
	req := new(dtos.CreateSupplierRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateCreateSupplier(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.SupplierService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	createdSupplier, err := h.SupplierService.CreateSupplier(req, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusCreated, "Supplier created successfully", createdSupplier)
}

func (h *SupplierHandler) UpdateSupplier(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid Uuid format")
	}
	req := new(dtos.UpdateSupplierRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}
	lang := c.Get("lang").(string)
	if messages := validators.ValidateUpdateSupplier(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	result, err := h.SupplierService.UpdateSupplier(parsedUuid, req, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "Supplier updated successfully", result)
}

func (h *SupplierHandler) DeleteSupplier(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid Uuid format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	err = h.SupplierService.DeleteSupplier(parsedUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "Supplier deleted successfully", nil)
}
