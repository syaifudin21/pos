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
	fieldToMessage := map[string]string{
		"Quantity": "quantity_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
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
	fieldToMessage := map[string]string{
		"OutletUuid":  "outlet_uuid_required",
		"Productuuid": "product_uuid_required",
		"Quantity":    "quantity_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
		}
	}
	return messages
}
