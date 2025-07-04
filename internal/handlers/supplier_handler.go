package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
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
// @Success 200 {object} SuccessResponse{data=[]models.Supplier}
// @Failure 500 {object} ErrorResponse
// @Router /suppliers [get]
func (h *SupplierHandler) GetAllSuppliers(c echo.Context) error {
	suppliers, err := h.SupplierService.GetAllSuppliers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Suppliers retrieved successfully", Data: suppliers})
}

// GetSupplierByuuid godoc
// @Summary Get supplier by Uuid
// @Description Get a single supplier by its Uuid.
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param uuid path string true "Supplier Uuid"
// @Success 200 {object} SuccessResponse{data=models.Supplier}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers/{uuid} [get]
func (h *SupplierHandler) GetSupplierByuuid(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Uuid format"})
	}
	supplier, err := h.SupplierService.GetSupplierByuuid(uuid)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Supplier retrieved successfully", Data: supplier})
}

// CreateSupplier godoc
// @Summary Create a new supplier
// @Description Create a new supplier with the provided details.
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param supplier body CreateSupplierRequest true "Supplier details"
// @Success 201 {object} SuccessResponse{data=models.Supplier}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers [post]
func (h *SupplierHandler) CreateSupplier(c echo.Context) error {
	req := new(CreateSupplierRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	newSupplier := &models.Supplier{
		Name:    req.Name,
		Contact: req.Contact,
		Address: req.Address,
	}

	createdSupplier, err := h.SupplierService.CreateSupplier(newSupplier)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, SuccessResponse{Message: "Supplier created successfully", Data: createdSupplier})
}

// UpdateSupplier godoc
// @Summary Update an existing supplier
// @Description Update an existing supplier by its Uuid.
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param uuid path string true "Supplier Uuid"
// @Param supplier body UpdateSupplierRequest true "Updated supplier details"
// @Success 200 {object} SuccessResponse{data=models.Supplier}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers/{uuid} [put]
func (h *SupplierHandler) UpdateSupplier(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Uuid format"})
	}
	req := new(UpdateSupplierRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	updatedSupplier := &models.Supplier{
		Name:    req.Name,
		Contact: req.Contact,
		Address: req.Address,
	}

	result, err := h.SupplierService.UpdateSupplier(uuid, updatedSupplier)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Supplier updated successfully", Data: result})
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
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Uuid format"})
	}
	err = h.SupplierService.DeleteSupplier(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNoContent, SuccessResponse{Message: "Supplier deleted successfully"})
}

type CreateSupplierRequest struct {
	Name    string `json:"name"`
	Contact string `json:"contact,omitempty"`
	Address string `json:"address,omitempty"`
}

type UpdateSupplierRequest struct {
	Name    string `json:"name"`
	Contact string `json:"contact,omitempty"`
	Address string `json:"address,omitempty"`
}
