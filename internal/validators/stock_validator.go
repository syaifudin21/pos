package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

var stockValidator = validator.New()

func ValidateUpdateStock(req *dtos.UpdateStockRequest) []string {
	// Custom validation logic
	if (req.ProductUuid == uuid.Nil && req.ProductVariantUuid == uuid.Nil) || (req.ProductUuid != uuid.Nil && req.ProductVariantUuid != uuid.Nil) {
		return []string{"either_product_uuid_or_product_variant_uuid_is_required"}
	}

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
			messages = append(messages, msg)
		}
	}
	return messages
}
