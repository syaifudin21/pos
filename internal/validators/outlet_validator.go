package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var outletValidator = validator.New()

func ValidateCreateOutlet(req *dtos.OutletCreateRequest, lang string) []string {
	err := outletValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Name":
			messages = append(messages, localization.GetLocalizedValidationMessage("name_required", lang))
		case "Address":
			messages = append(messages, localization.GetLocalizedValidationMessage("address_required", lang))
		case "Type":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_type_required", lang))
		}
	}
	return messages
}

func ValidateUpdateOutlet(req *dtos.OutletUpdateRequest, lang string) []string {
	err := outletValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Name":
			messages = append(messages, localization.GetLocalizedValidationMessage("name_required", lang))
		case "Address":
			messages = append(messages, localization.GetLocalizedValidationMessage("address_required", lang))
		case "Type":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_type_required", lang))
		}
	}
	return messages
}
