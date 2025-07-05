package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var orderValidator = validator.New()

func ValidateCreateOrder(req *dtos.CreateOrderRequest, lang string) []string {
	err := orderValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "OutletUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("outlet_uuid_required", lang))
		case "Items":
			messages = append(messages, localization.GetLocalizedValidationMessage("order_items_required", lang))
		case "PaymentMethod":
			messages = append(messages, localization.GetLocalizedValidationMessage("payment_method_required", lang))
		}
	}
	return messages
}

func ValidateOrderItem(req *dtos.OrderItemRequest, lang string) []string {
	err := orderValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "ProductUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_uuid_required", lang))
		case "Quantity":
			messages = append(messages, localization.GetLocalizedValidationMessage("quantity_required", lang))
		}
	}
	return messages
}
