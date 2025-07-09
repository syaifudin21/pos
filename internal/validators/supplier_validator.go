package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

var supplierValidator = validator.New()

func ValidateCreateSupplier(req *dtos.CreateSupplierRequest) []string {
	err := supplierValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"Name":    "name_required",
		"Contact": "contact_required",
		"Address": "address_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}

func ValidateUpdateSupplier(req *dtos.UpdateSupplierRequest) []string {
	err := supplierValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"Name":    "name_required",
		"Contact": "contact_required",
		"Address": "address_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		}
	}
	return messages
}