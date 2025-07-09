package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

var productValidator = validator.New()

func ValidateCreateProduct(req *dtos.ProductCreateRequest) []string {
	err := productValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"Name":        "product_name_required",
		"Description": "product_description_required",
		"Price":       "product_price_required",
		"SKU":         "product_sku_required",
		"Type":        "product_type_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}

func ValidateUpdateProduct(req *dtos.ProductUpdateRequest) []string {
	err := productValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"Name":        "product_name_required",
		"Description": "product_description_required",
		"Price":       "product_price_required",
		"SKU":         "product_sku_required",
		"Type":        "product_type_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}