package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
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
// @Success 200 {object} SuccessResponse{data=[]dtos.OutletResponse}
// @Failure 500 {object} ErrorResponse
// @Router /outlets [get]
func (h *OutletHandler) GetAllOutlets(c echo.Context) error {
	outlets, err := h.OutletService.GetAllOutlets()
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	var outletResponses []dtos.OutletResponse
	for _, outlet := range outlets {
		outletResponses = append(outletResponses, dtos.OutletResponse{
			ID:      outlet.ID,
			Uuid:    outlet.Uuid,
			Name:    outlet.Name,
			Address: outlet.Address,
			Type:    outlet.Type,
		})
	}
	return JSONSuccess(c, http.StatusOK, "outlets_retrieved_successfully", outletResponses)
}

// GetOutletByID godoc
// @Summary Get outlet by Uuid
// @Description Get a single outlet by its Uuid.
// @Tags Outlets
// @Accept json
// @Produce json
// @Param uuid path string true "Outlet Uuid"
// @Success 200 {object} SuccessResponse{data=dtos.OutletResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{uuid} [get]
func (h *OutletHandler) GetOutletByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	outlet, err := h.OutletService.GetOutletByUuid(id)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "outlet_retrieved_successfully", outlet)
}

// CreateOutlet godoc
// @Summary Create a new outlet
// @Description Create a new outlet with the provided details.
// @Tags Outlets
// @Accept json
// @Produce json
// @Param outlet body dtos.OutletCreateRequest true "Outlet details"
// @Success 201 {object} SuccessResponse{data=dtos.OutletResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets [post]
func (h *OutletHandler) CreateOutlet(c echo.Context) error {
	outlet := new(dtos.OutletCreateRequest)
	if err := c.Bind(outlet); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	createdOutlet, err := h.OutletService.CreateOutlet(outlet)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusCreated, "outlet_created_successfully", dtos.OutletResponse{
		ID:      createdOutlet.ID,
		Uuid:    createdOutlet.Uuid,
		Name:    createdOutlet.Name,
		Address: createdOutlet.Address,
		Type:    createdOutlet.Type,
	})
}

// UpdateOutlet godoc
// @Summary Update an existing outlet
// @Description Update an existing outlet by its Uuid.
// @Tags Outlets
// @Accept json
// @Produce json
// @Param uuid path string true "Outlet Uuid"
// @Param outlet body dtos.OutletUpdateRequest true "Updated outlet details"
// @Success 200 {object} SuccessResponse{data=dtos.OutletResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{uuid} [put]
func (h *OutletHandler) UpdateOutlet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	outlet := new(dtos.OutletUpdateRequest)
	if err := c.Bind(outlet); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	result, err := h.OutletService.UpdateOutlet(id, outlet)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "outlet_updated_successfully", result)
}

// DeleteOutlet godoc
// @Summary Delete an outlet
// @Description Delete an outlet by its Uuid.
// @Tags Outlets
// @Accept json
// @Produce json
// @Param uuid path string true "Outlet Uuid"
// @Success 204 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{uuid} [delete]
func (h *OutletHandler) DeleteOutlet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	err = h.OutletService.DeleteOutlet(id)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "outlet_deleted_successfully", nil)
}
