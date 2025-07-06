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

type OutletHandler struct {
	OutletService *services.OutletService
}

func NewOutletHandler(outletService *services.OutletService) *OutletHandler {
	return &OutletHandler{OutletService: outletService}
}

func (h *OutletHandler) GetAllOutlets(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.OutletService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	outlets, err := h.OutletService.GetAllOutlets(ownerID)
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

func (h *OutletHandler) GetOutletByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.OutletService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	outlet, err := h.OutletService.GetOutletByUuid(id, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "outlet_retrieved_successfully", outlet)
}

func (h *OutletHandler) CreateOutlet(c echo.Context) error {
	outlet := new(dtos.OutletCreateRequest)
	if err := c.Bind(outlet); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateCreateOutlet(outlet, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.OutletService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	createdOutlet, err := h.OutletService.CreateOutlet(outlet, ownerID)
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

func (h *OutletHandler) UpdateOutlet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	outlet := new(dtos.OutletUpdateRequest)
	if err := c.Bind(outlet); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateUpdateOutlet(outlet, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.OutletService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	result, err := h.OutletService.UpdateOutlet(id, outlet, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "outlet_updated_successfully", result)
}

func (h *OutletHandler) DeleteOutlet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	err = h.OutletService.DeleteOutlet(id, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "outlet_deleted_successfully", nil)
}
