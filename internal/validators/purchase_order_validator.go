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
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "SupplierUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("supplier_uuid_required", lang))
		case "OutletUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("outlet_uuid_required", lang))
		case "Items":
			messages = append(messages, localization.GetLocalizedValidationMessage("purchase_items_required", lang))
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
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "ProductUuid":
			messages = append(messages, localization.GetLocalizedValidationMessage("product_uuid_required", lang))
		case "Quantity":
			messages = append(messages, localization.GetLocalizedValidationMessage("quantity_required", lang))
		case "Price":
			messages = append(messages, localization.GetLocalizedValidationMessage("price_required", lang))
		}
	}
	return messages
}
