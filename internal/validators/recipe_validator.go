package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

var recipeValidator = validator.New()

func ValidateCreateRecipe(req *dtos.CreateRecipeRequest) []string {
	err := recipeValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"MainProductUuid": "main_product_uuid_required",
		"ComponentUuid":   "component_uuid_required",
		"Quantity":        "quantity_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}

func ValidateUpdateRecipe(req *dtos.UpdateRecipeRequest) []string {
	err := recipeValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"MainProductUuid": "main_product_uuid_required",
		"ComponentUuid":   "component_uuid_required",
		"Quantity":        "quantity_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}