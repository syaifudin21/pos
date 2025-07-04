package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/services"
)

type OutletHandler struct {
	OutletService *services.OutletService
}

func NewOutletHandler(outletService *services.OutletService) *OutletHandler {
	return &OutletHandler{OutletService: outletService}
}

// GetAllOutlets godoc
// @Summary Get all outlets
// @Description Get a list of all outlets.
// @Tags Outlets
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse{data=[]models.Outlet}
// @Failure 500 {object} ErrorResponse
// @Router /outlets [get]
func (h *OutletHandler) GetAllOutlets(c echo.Context) error {
	outlets, err := h.OutletService.GetAllOutlets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Outlets retrieved successfully", Data: outlets})
}

// GetOutletByID godoc
// @Summary Get outlet by External ID
// @Description Get a single outlet by its External ID.
// @Tags Outlets
// @Accept json
// @Produce json
// @Param uuid path string true "Outlet External ID (UUID)"
// @Success 200 {object} SuccessResponse{data=models.Outlet}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{uuid} [get]
func (h *OutletHandler) GetOutletByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid UUID format"})
	}
	outlet, err := h.OutletService.GetOutletByUuid(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Outlet retrieved successfully", Data: outlet})
}

// CreateOutlet godoc
// @Summary Create a new outlet
// @Description Create a new outlet with the provided details.
// @Tags Outlets
// @Accept json
// @Produce json
// @Param outlet body OutletCreateRequest true "Outlet details"
// @Success 201 {object} SuccessResponse{data=models.Outlet}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets [post]
func (h *OutletHandler) CreateOutlet(c echo.Context) error {
	outlet := new(OutletCreateRequest)
	if err := c.Bind(outlet); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	newOutlet := &models.Outlet{
		Name:    outlet.Name,
		Address: outlet.Address,
		Type:    outlet.Type,
	}

	createdOutlet, err := h.OutletService.CreateOutlet(newOutlet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, SuccessResponse{Message: "Outlet created successfully", Data: createdOutlet})
}

// UpdateOutlet godoc
// @Summary Update an existing outlet
// @Description Update an existing outlet by its External ID.
// @Tags Outlets
// @Accept json
// @Produce json
// @Param uuid path string true "Outlet External ID (UUID)"
// @Param outlet body OutletUpdateRequest true "Updated outlet details"
// @Success 200 {object} SuccessResponse{data=models.Outlet}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{uuid} [put]
func (h *OutletHandler) UpdateOutlet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid UUID format"})
	}
	outlet := new(OutletUpdateRequest)
	if err := c.Bind(outlet); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	updatedOutlet := &models.Outlet{
		Name:    outlet.Name,
		Address: outlet.Address,
		Type:    outlet.Type,
	}

	result, err := h.OutletService.UpdateOutlet(id, updatedOutlet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Outlet updated successfully", Data: result})
}

// DeleteOutlet godoc
// @Summary Delete an outlet
// @Description Delete an outlet by its External ID.
// @Tags Outlets
// @Accept json
// @Produce json
// @Param uuid path string true "Outlet External ID (UUID)"
// @Success 204 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{uuid} [delete]
func (h *OutletHandler) DeleteOutlet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid UUID format"})
	}
	err = h.OutletService.DeleteOutlet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNoContent, SuccessResponse{Message: "Outlet deleted successfully"})
}

type OutletCreateRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type OutletUpdateRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Type    string `json:"type"`
}
