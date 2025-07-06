package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var purchaseOrderValidator = validator.New()

func ValidateCreatePurchaseOrder(req *dtos.CreatePurchaseOrderRequest, lang string) []string {
	err := purchaseOrderValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"SupplierUuid": "supplier_uuid_required",
		"OutletUuid":   "outlet_uuid_required",
		"Items":        "purchase_items_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
		}
	}
	return messages
}

func ValidatePurchaseItem(req *dtos.PurchaseItemRequest, lang string) []string {
	err := purchaseOrderValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"ProductUuid": "product_uuid_required",
		"Quantity":    "quantity_required",
		"Price":       "price_required",
	}
	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
		}
	}
	return messages
}
