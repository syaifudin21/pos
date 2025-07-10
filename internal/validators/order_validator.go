package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
		"OutletUuid":      "outlet_uuid_required",
		"Items":           "order_items_required",
		"PaymentMethodID": "payment_method_id_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}

func ValidateOrderItem(req *dtos.OrderItemRequest) []string {
	// Custom validation logic
	if (req.ProductUuid == uuid.Nil && req.ProductVariantUuid == uuid.Nil) || (req.ProductUuid != uuid.Nil && req.ProductVariantUuid != uuid.Nil) {
		return []string{"either_product_uuid_or_product_variant_uuid_is_required_for_order_item"}
	}

	err := orderValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"Quantity":    "quantity_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}