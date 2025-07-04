package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
// @Success 200 {object} SuccessResponse{data=models.Recipe}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /recipes/{uuid} [get]
func (h *RecipeHandler) GetRecipeByUuid(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid UUID format"})
	}
	recipe, err := h.RecipeService.GetRecipeByUuid(uuid)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Recipe retrieved successfully", Data: recipe})
}

// GetRecipesByMainProduct godoc
// @Summary Get recipes by main product
// @Description Get a list of recipes for a given main product.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param main_product_uuid path string true "Main Product Uuid"
// @Success 200 {object} SuccessResponse{data=[]models.Recipe}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{main_product_uuid}/recipes [get]
func (h *RecipeHandler) GetRecipesByMainProduct(c echo.Context) error {
	mainProductUuid, err := uuid.Parse(c.Param("main_product_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Main Product Uuid format"})
	}
	recipes, err := h.RecipeService.GetRecipesByMainProduct(mainProductUuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Recipes retrieved successfully", Data: recipes})
}

// CreateRecipe godoc
// @Summary Create a new recipe
// @Description Create a new recipe with the provided details.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param recipe body CreateRecipeRequest true "Recipe details"
// @Success 201 {object} SuccessResponse{data=models.Recipe}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /recipes [post]
func (h *RecipeHandler) CreateRecipe(c echo.Context) error {
	req := new(CreateRecipeRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	createdRecipe, err := h.RecipeService.CreateRecipe(req.MainProductUuid, req.ComponentUuid, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, SuccessResponse{Message: "Recipe created successfully", Data: createdRecipe})
}

// UpdateRecipe godoc
// @Summary Update an existing recipe
// @Description Update an existing recipe by its Uuid.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param uuid path string true "Recipe Uuid"
// @Param recipe body UpdateRecipeRequest true "Updated recipe details"
// @Success 200 {object} SuccessResponse{data=models.Recipe}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /recipes/{uuid} [put]
func (h *RecipeHandler) UpdateRecipe(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid UUID format"})
	}
	req := new(UpdateRecipeRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	updatedRecipe, err := h.RecipeService.UpdateRecipe(uuid, req.MainProductUuid, req.ComponentUuid, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Recipe updated successfully", Data: updatedRecipe})
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
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid UUID format"})
	}
	err = h.RecipeService.DeleteRecipe(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNoContent, SuccessResponse{Message: "Recipe deleted successfully"})
}

type CreateRecipeRequest struct {
	MainProductUuid uuid.UUID `json:"main_product_uuid"`
	ComponentUuid   uuid.UUID `json:"component_uuid"`
	Quantity        float64   `json:"quantity"`
}

type UpdateRecipeRequest struct {
	MainProductUuid uuid.UUID `json:"main_product_uuid"`
	ComponentUuid   uuid.UUID `json:"component_uuid"`
	Quantity        float64   `json:"quantity"`
}
