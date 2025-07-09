package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type RecipeHandler struct {
	RecipeService      *services.RecipeService
	UserContextService *services.UserContextService
}

func NewRecipeHandler(recipeService *services.RecipeService, userContextService *services.UserContextService) *RecipeHandler {
	return &RecipeHandler{RecipeService: recipeService, UserContextService: userContextService}
}
func (h *RecipeHandler) GetRecipeByUuid(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	recipe, err := h.RecipeService.GetRecipeByUuid(parsedUuid, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "recipe_retrieved_successfully", recipe)
}

func (h *RecipeHandler) GetRecipesByMainProduct(c echo.Context) error {
	mainProductUuidParam := c.Param("main_product_uuid")
	mainProductUuid, err := uuid.Parse(mainProductUuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_main_product_uuid_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	recipes, err := h.RecipeService.GetRecipesByMainProduct(mainProductUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	// The service now returns []dtos.RecipeResponse directly, so no need for manual mapping here.
	return JSONSuccess(c, http.StatusOK, "recipes_retrieved_successfully", recipes)
}
func (h *RecipeHandler) CreateRecipe(c echo.Context) error {
	req, ok := c.Get("validated_data").(*dtos.CreateRecipeRequest)
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

	createdRecipe, err := h.RecipeService.CreateRecipe(req.MainProductUuid, req.ComponentUuid, req.Quantity, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusCreated, "recipe_created_successfully", createdRecipe)
}
func (h *RecipeHandler) UpdateRecipe(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	req, ok := c.Get("validated_data").(*dtos.UpdateRecipeRequest)
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

	updatedRecipe, err := h.RecipeService.UpdateRecipe(parsedUuid, req.MainProductUuid, req.ComponentUuid, req.Quantity, ownerID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "recipe_updated_successfully", updatedRecipe)
}

func (h *RecipeHandler) DeleteRecipe(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	err = h.RecipeService.DeleteRecipe(parsedUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "recipe_deleted_successfully", nil)
}
