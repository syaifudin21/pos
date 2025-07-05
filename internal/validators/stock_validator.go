package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var stockValidator = validator.New()

func ValidateUpdateStock(req *dtos.UpdateStockRequest, lang string) []string {
	err := stockValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Quantity":
			messages = append(messages, localization.GetLocalizedValidationMessage("quantity_required", lang))
		}
	}
	return messages
}

func ValidateGlobalStockUpdate(req *dtos.GlobalStockUpdateRequest, lang string) []string {
	err := stockValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "OutletUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("outlet_uuid_required", lang))
		case "Productuuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_uuid_required", lang))
		case "Quantity":
			messages = append(messages, localization.GetLocalizedValidationMessage("quantity_required", lang))
		}
	}
	return messages
}
