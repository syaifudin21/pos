package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type SupplierHandler struct {
	SupplierService *services.SupplierService
}

func NewSupplierHandler(supplierService *services.SupplierService) *SupplierHandler {
	return &SupplierHandler{SupplierService: supplierService}
}

// GetAllSuppliers godoc
// @Summary Get all suppliers
// @Description Get a list of all suppliers.
// @Tags Suppliers
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse{data=[]dtos.SupplierResponse}
// @Failure 500 {object} ErrorResponse
// @Router /suppliers [get]
func (h *SupplierHandler) GetAllSuppliers(c echo.Context) error {
	suppliers, err := h.SupplierService.GetAllSuppliers()
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

// GetSupplierByuuid godoc
// @Summary Get supplier by Uuid
// @Description Get a single supplier by its Uuid.
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param uuid path string true "Supplier Uuid"
// @Success 200 {object} SuccessResponse{data=dtos.SupplierResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers/{uuid} [get]
func (h *SupplierHandler) GetSupplierByuuid(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid Uuid format")
	}
	supplier, err := h.SupplierService.GetSupplierByuuid(uuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "Supplier retrieved successfully", supplier)
}

// CreateSupplier godoc
// @Summary Create a new supplier
// @Description Create a new supplier with the provided details.
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param supplier body dtos.CreateSupplierRequest true "Supplier details"
// @Success 201 {object} SuccessResponse{data=dtos.SupplierResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers [post]
func (h *SupplierHandler) CreateSupplier(c echo.Context) error {
	req := new(dtos.CreateSupplierRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid request payload")
	}

	createdSupplier, err := h.SupplierService.CreateSupplier(req)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusCreated, "Supplier created successfully", createdSupplier)
}

// UpdateSupplier godoc
// @Summary Update an existing supplier
// @Description Update an existing supplier by its Uuid.
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param uuid path string true "Supplier Uuid"
// @Param supplier body dtos.UpdateSupplierRequest true "Updated supplier details"
// @Success 200 {object} SuccessResponse{data=dtos.SupplierResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers/{uuid} [put]
func (h *SupplierHandler) UpdateSupplier(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid Uuid format")
	}
	req := new(dtos.UpdateSupplierRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid request payload")
	}

	result, err := h.SupplierService.UpdateSupplier(uuid, req)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "Supplier updated successfully", result)
}

// DeleteSupplier godoc
// @Summary Delete a supplier
// @Description Delete a supplier by its Uuid.
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param uuid path string true "Supplier Uuid"
// @Success 204 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers/{uuid} [delete]
func (h *SupplierHandler) DeleteSupplier(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid Uuid format")
	}
	err = h.SupplierService.DeleteSupplier(uuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "Supplier deleted successfully", nil)
}