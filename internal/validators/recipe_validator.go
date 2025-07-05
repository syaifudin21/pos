package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var recipeValidator = validator.New()

func ValidateCreateRecipe(req *dtos.CreateRecipeRequest, lang string) []string {
	err := recipeValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "MainProductUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("main_product_uuid_required", lang))
		case "ComponentUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("component_uuid_required", lang))
		case "Quantity":
			messages = append(messages, localization.GetLocalizedValidationMessage("quantity_required", lang))
		}
	}
	return messages
}

func ValidateUpdateRecipe(req *dtos.UpdateRecipeRequest, lang string) []string {
	err := recipeValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "MainProductUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("main_product_uuid_required", lang))
		case "ComponentUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("component_uuid_required", lang))
		case "Quantity":
			messages = append(messages, localization.GetLocalizedValidationMessage("quantity_required", lang))
		}
	}
	return messages
}
