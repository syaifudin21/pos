package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var productValidator = validator.New()

func ValidateCreateProduct(req *dtos.ProductCreateRequest, lang string) []string {
	err := productValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Name":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_name_required", lang))
		case "Description":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_description_required", lang))
		case "Price":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_price_required", lang))
		case "SKU":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_sku_required", lang))
		case "Type":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_type_required", lang))
		}
	}
	return messages
}

func ValidateUpdateProduct(req *dtos.ProductUpdateRequest, lang string) []string {
	err := productValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Name":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_name_required", lang))
		case "Description":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_description_required", lang))
		case "Price":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_price_required", lang))
		case "SKU":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_sku_required", lang))
		case "Type":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_type_required", lang))
		}
	}
	return messages
}
