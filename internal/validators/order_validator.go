package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

var orderValidator = validator.New()

func ValidateCreateOrder(req *dtos.CreateOrderRequest) []string {
	err := orderValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"OutletUuid":    "outlet_uuid_required",
		"Items":         "order_items_required",
		"PaymentMethod": "payment_method_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}

func ValidateOrderItem(req *dtos.OrderItemRequest) []string {
	err := orderValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"ProductUuid": "product_uuid_required",
		"Quantity":    "quantity_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}