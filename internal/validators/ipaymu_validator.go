package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var ipaymuValidator = validator.New()

func ValidateIpaymuRequest(req *dtos.IpaymuRequest, lang string) []string {
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
		case "ReturnUrl":
			messages = append(messages, localization.GetLocalizedValidationMessage("return_url_required", lang))
		case "CancelUrl":
			messages = append(messages, localization.GetLocalizedValidationMessage("cancel_url_required", lang))
		case "NotifyUrl":
			messages = append(messages, localization.GetLocalizedValidationMessage("notify_url_required", lang))
		case "ReferenceId":
			messages = append(messages, localization.GetLocalizedValidationMessage("reference_id_required", lang))
		case "BuyerName":
			messages = append(messages, localization.GetLocalizedValidationMessage("buyer_name_required", lang))
		case "BuyerEmail":
			messages = append(messages, localization.GetLocalizedValidationMessage("buyer_email_required", lang))
		case "BuyerPhone":
			messages = append(messages, localization.GetLocalizedValidationMessage("buyer_phone_required", lang))
		case "Udf1":
			messages = append(messages, localization.GetLocalizedValidationMessage("udf1_required", lang))
		}
	}
	return messages
}
