package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var supplierValidator = validator.New()

func ValidateCreateSupplier(req *dtos.CreateSupplierRequest, lang string) []string {
	err := supplierValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Name":
			messages = append(messages, localization.GetLocalizedValidationMessage("name_required", lang))
		case "Contact":
			messages = append(messages, localization.GetLocalizedValidationMessage("contact_required", lang))
		case "Address":
			messages = append(messages, localization.GetLocalizedValidationMessage("address_required", lang))
		}
	}
	return messages
}

func ValidateUpdateSupplier(req *dtos.UpdateSupplierRequest, lang string) []string {
	err := supplierValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Name":
			messages = append(messages, localization.GetLocalizedValidationMessage("name_required", lang))
		case "Contact":
			messages = append(messages, localization.GetLocalizedValidationMessage("contact_required", lang))
		case "Address":
			messages = append(messages, localization.GetLocalizedValidationMessage("address_required", lang))
		}
	}
	return messages
}
