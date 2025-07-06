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
	tagToMessage := map[string]func(field string) string{
		"required": func(field string) string { return localization.GetLocalizedValidationMessage(field+"_required", lang) },
		"email":    func(_ string) string { return localization.GetLocalizedValidationMessage("email_invalid", lang) },
		"uuid":     func(_ string) string { return localization.GetLocalizedValidationMessage("uuid_invalid", lang) },
		"url":      func(_ string) string { return localization.GetLocalizedValidationMessage("url_invalid", lang) },
		"gt": func(field string) string {
			return localization.GetLocalizedValidationMessage(field+"_greater_than_zero", lang)
		},
		"dive": func(field string) string {
			return localization.GetLocalizedValidationMessage(field+"_dive_required", lang)
		},
		"required_if": func(field string) string {
			return localization.GetLocalizedValidationMessage(field+"_required_if", lang)
		},
		"required_with": func(field string) string {
			return localization.GetLocalizedValidationMessage(field+"_required_with", lang)
		},
	}
	minFieldToMessage := map[string]string{
		"Product": "product_min_one",
		"Qty":     "qty_min_one",
		"Price":   "price_min_one",
	}
	fieldToMessage := map[string]string{
		"Product":  "product_required",
		"Qty":      "qty_required",
		"Price":    "price_required",
		"Name":     "name_required",
		"Email":    "email_required",
		"Phone":    "phone_required",
		"Callback": "callback_required",
		"Method":   "method_required",
		"Channel":  "channel_required",
	}

	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
		}
		if err.Tag() == "min" {
			if msg, ok := minFieldToMessage[err.Field()]; ok {
				messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
			}
		} else if fn, ok := tagToMessage[err.Tag()]; ok {
			messages = append(messages, fn(err.Field()))
		}
	}
	return messages
}

func ValidateIpaymuNotify(req *dtos.IpaymuNotifyRequest, lang string) []string {
	err := ipaymuValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	tagToMessage := map[string]func(field string) string{
		"required": func(field string) string { return localization.GetLocalizedValidationMessage(field+"_required", lang) },
		"uuid":     func(_ string) string { return localization.GetLocalizedValidationMessage("uuid_invalid", lang) },
	}

	fieldToMessage := map[string]string{
		"ReferenceId": "reference_id_required",
	}

	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
		} else if fn, ok := tagToMessage[err.Tag()]; ok {
			messages = append(messages, fn(err.Field()))
		}
	}
	return messages
}
