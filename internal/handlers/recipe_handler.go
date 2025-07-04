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

type RecipeHandler struct {
	RecipeService *services.RecipeService
}

func NewRecipeHandler(recipeService *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{RecipeService: recipeService}
}
func (h *RecipeHandler) GetRecipeByUuid(c echo.Context) error {
	uuidParam := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.RecipeService.GetOwnerID(userID)
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.RecipeService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	recipes, err := h.RecipeService.GetRecipesByMainProduct(mainProductUuid, ownerID)
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
func (h *RecipeHandler) CreateRecipe(c echo.Context) error {
	req := new(dtos.CreateRecipeRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateCreateRecipe(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.RecipeService.GetOwnerID(userID)
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
	req := new(dtos.UpdateRecipeRequest)
	if err := c.Bind(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	lang := c.Get("lang").(string)
	if messages := validators.ValidateUpdateRecipe(req, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.RecipeService.GetOwnerID(userID)
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	err = h.RecipeService.DeleteRecipe(parsedUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "recipe_deleted_successfully", nil)
}
