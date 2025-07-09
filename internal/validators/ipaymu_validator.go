package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

var ipaymuValidator = validator.New()

func ValidateCreateDirectPayment(req *dtos.CreateDirectPaymentRequest) []string {
	err := ipaymuValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	tagToMessage := map[string]func(field string) string{
		"required": func(field string) string { return field+"_required" },
		"email":    func(_ string) string { return "email_invalid" },
		"uuid":     func(_ string) string { return "uuid_invalid" },
		"url":      func(_ string) string { return "url_invalid" },
		"gt": func(field string) string {
			return field+"_greater_than_zero"
		},
		"dive": func(field string) string {
			return field+"_dive_required"
		},
		"required_if": func(field string) string {
			return field+"_required_if"
		},
		"required_with": func(field string) string {
			return field+"_required_with"
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
			messages = append(messages, msg)
		}
		if err.Tag() == "min" {
			if msg, ok := minFieldToMessage[err.Field()]; ok {
				messages = append(messages, msg)
			}
		} else if fn, ok := tagToMessage[err.Tag()]; ok {
			messages = append(messages, fn(err.Field()))
		}
	}
	return messages
}

func ValidateIpaymuNotify(req *dtos.IpaymuNotifyRequest) []string {
	err := ipaymuValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	tagToMessage := map[string]func(field string) string{
		"required": func(field string) string { return field+"_required" },
		"uuid":     func(_ string) string { return "uuid_invalid" },
	}

	fieldToMessage := map[string]string{
		"ReferenceId": "reference_id_required",
	}

	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		} else if fn, ok := tagToMessage[err.Tag()]; ok {
			messages = append(messages, fn(err.Field()))
		}
	}
	return messages
}

func ValidateRegisterIpaymu(req *dtos.RegisterIpaymuRequest) []string {
	err := ipaymuValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	fieldToMessage := map[string]string{
		"Name":         "name_required",
		"Phone":        "phone_required",
		"Password":     "password_required",
		"WithoutEmail": "without_email_required",
	}
	tagToMessage := map[string]func(field string) string{
		"required": func(field string) string { return field+"_required" },
		"email":    func(_ string) string { return "email_invalid" },
		"oneof":    func(field string) string { return field+"_oneof_invalid" },
	}

	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			messages = append(messages, msg)
		} else if fn, ok := tagToMessage[err.Tag()]; ok {
			messages = append(messages, fn(err.Field()))
		}
	}
	return messages
}
