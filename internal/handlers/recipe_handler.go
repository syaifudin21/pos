package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
)

type RecipeHandler struct {
	RecipeService *services.RecipeService
}

func NewRecipeHandler(recipeService *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{RecipeService: recipeService}
}

// GetRecipeByUuid godoc
// @Summary Get recipe by Uuid
// @Description Get a single recipe by its Uuid.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param uuid path string true "Recipe Uuid"
// @Success 200 {object} SuccessResponse{data=dtos.RecipeResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /recipes/{uuid} [get]
func (h *RecipeHandler) GetRecipeByUuid(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	recipe, err := h.RecipeService.GetRecipeByUuid(uuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "recipe_retrieved_successfully", recipe)
}

// GetRecipesByMainProduct godoc
// @Summary Get recipes by main product
// @Description Get a list of recipes for a given main product.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param main_product_uuid path string true "Main Product Uuid"
// @Success 200 {object} SuccessResponse{data=[]dtos.RecipeResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{main_product_uuid}/recipes [get]
func (h *RecipeHandler) GetRecipesByMainProduct(c echo.Context) error {
	mainProductUuid, err := uuid.Parse(c.Param("main_product_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_main_product_uuid_format")
	}
	recipes, err := h.RecipeService.GetRecipesByMainProduct(mainProductUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	var recipeResponses []dtos.RecipeResponse
	for _, recipe := range recipes {
		recipeResponses = append(recipeResponses, dtos.RecipeResponse{
			ID:              recipe.ID,
			Uuid:            recipe.Uuid,
			MainProductID:   recipe.MainProductID,
			MainProductUuid: recipe.MainProductUuid,
			ComponentID:     recipe.ComponentID,
			ComponentUuid:   recipe.ComponentUuid,
			Quantity:        recipe.Quantity,
		})
	}
	return JSONSuccess(c, http.StatusOK, "recipes_retrieved_successfully", recipeResponses)
}

// CreateRecipe godoc
// @Summary Create a new recipe
// @Description Create a new recipe with the provided details.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param recipe body dtos.CreateRecipeRequest true "Recipe details"
// @Success 201 {object} SuccessResponse{data=dtos.RecipeResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /recipes [post]
func (h *RecipeHandler) CreateRecipe(c echo.Context) error {
	req := new(dtos.CreateRecipeRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	createdRecipe, err := h.RecipeService.CreateRecipe(req.MainProductUuid, req.ComponentUuid, req.Quantity)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusCreated, "recipe_created_successfully", createdRecipe)
}

// UpdateRecipe godoc
// @Summary Update an existing recipe
// @Description Update an existing recipe by its Uuid.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param uuid path string true "Recipe Uuid"
// @Param recipe body dtos.UpdateRecipeRequest true "Updated recipe details"
// @Success 200 {object} SuccessResponse{data=dtos.RecipeResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /recipes/{uuid} [put]
func (h *RecipeHandler) UpdateRecipe(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	req := new(dtos.UpdateRecipeRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	updatedRecipe, err := h.RecipeService.UpdateRecipe(uuid, req.MainProductUuid, req.ComponentUuid, req.Quantity)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "recipe_updated_successfully", updatedRecipe)
}

// DeleteRecipe godoc
// @Summary Delete a recipe
// @Description Delete a recipe by its Uuid.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param uuid path string true "Recipe Uuid"
// @Success 204 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /recipes/{uuid} [delete]
func (h *RecipeHandler) DeleteRecipe(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	err = h.RecipeService.DeleteRecipe(uuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "recipe_deleted_successfully", nil)
}
