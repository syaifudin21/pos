package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type SupplierHandler struct {
	SupplierService    *services.SupplierService
	UserContextService *services.UserContextService
}

func NewSupplierHandler(supplierService *services.SupplierService, userContextService *services.UserContextService) *SupplierHandler {
	return &SupplierHandler{SupplierService: supplierService, UserContextService: userContextService}
}

func (h *SupplierHandler) GetAllSuppliers(c echo.Context) error {
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
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

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
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
	req, ok := c.Get("validated_data").(*dtos.CreateSupplierRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
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
	req, ok := c.Get("validated_data").(*dtos.UpdateSupplierRequest)
	if !ok {
		return JSONError(c, http.StatusInternalServerError, "failed_to_get_validated_request")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

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

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	err = h.SupplierService.DeleteSupplier(parsedUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "Supplier deleted successfully", nil)
}
