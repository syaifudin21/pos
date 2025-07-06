package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var ipaymuValidator = validator.New()

func ValidateCreateDirectPayment(req *dtos.CreateDirectPaymentRequest, lang string) []string {
	err := ipaymuValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Product":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_required", lang))
		case "Qty":
			messages = append(messages, localization.GetLocalizedValidationMessage("qty_required", lang))
		case "Price":
			messages = append(messages, localization.GetLocalizedValidationMessage("price_required", lang))
		case "Name":
			messages = append(messages, localization.GetLocalizedValidationMessage("name_required", lang))
		case "Email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_required", lang))
		case "Phone":
			messages = append(messages, localization.GetLocalizedValidationMessage("phone_required", lang))
		case "Callback":
			messages = append(messages, localization.GetLocalizedValidationMessage("callback_required", lang))
		case "Method":
			messages = append(messages, localization.GetLocalizedValidationMessage("method_required", lang))
		case "Channel":
			messages = append(messages, localization.GetLocalizedValidationMessage("channel_required", lang))
		}
	}
	return messages
}
