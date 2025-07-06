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
	fieldToMessage := map[string]string{
		"Name":    "name_required",
		"Address": "address_required",
		"Type":    "product_type_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
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
	fieldToMessage := map[string]string{
		"Name":    "name_required",
		"Address": "address_required",
		"Type":    "product_type_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
		}
	}
	return messages
}
