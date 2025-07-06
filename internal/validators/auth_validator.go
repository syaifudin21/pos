package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/localization"
)

var authValidator = validator.New()

func init() {
	authValidator.RegisterValidation("passwordstrength", isPasswordStrong)
}

func isPasswordStrong(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Must contain at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	// Must contain at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	// Must contain at least one digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}

	// Must contain at least one special character
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(password) {
		return false
	}

	return true
}

func ValidateRegisterRequest(req *dtos.RegisterRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	tagToMessage := map[string]func(field string) string{
		"required":         func(field string) string { return localization.GetLocalizedValidationMessage(field+"_required", lang) },
		"passwordstrength": func(_ string) string { return localization.GetLocalizedValidationMessage("password_strength", lang) },
	}
	fieldToMessage := map[string]string{
		"Username": "username_required",
		"Password": "password_required",
		"Role":     "role_required",
	}

	for _, err := range err.(validator.ValidationErrors) {
		if msg, ok := fieldToMessage[err.Field()]; ok {
			// Special case for password strength
			if err.Field() == "Password" && err.Tag() == "passwordstrength" {
				messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
				continue
			}
			messages = append(messages, localization.GetLocalizedValidationMessage(msg, lang))
		} else if fn, ok := tagToMessage[err.Tag()]; ok {
			messages = append(messages, fn(err.Field()))
		}
	}
	return messages
}

func ValidateLoginRequest(req *dtos.LoginRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Username":
			messages = append(messages, localization.GetLocalizedValidationMessage("username_required", lang))
		case "Password":
			if err.Tag() == "passwordstrength" {
				messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
			} else {
				messages = append(messages, localization.GetLocalizedValidationMessage("password_required", lang))
			}
		}
	}
	return messages
}

func ValidateUpdateUserRequest(req *dtos.UpdateUserRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Username":
			messages = append(messages, localization.GetLocalizedValidationMessage("username_required", lang))
		case "Password":
			if err.Tag() == "passwordstrength" {
				messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
			} else {
				messages = append(messages, localization.GetLocalizedValidationMessage("password_required", lang))
			}
		case "Role":
			messages = append(messages, localization.GetLocalizedValidationMessage("role_invalid", lang))
		case "OutletID":
			messages = append(messages, localization.GetLocalizedValidationMessage("outlet_id_required", lang))
		}
	}
	return messages
}

func ValidateRegisterAdminRequest(req *dtos.RegisterAdminRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Username":
			messages = append(messages, localization.GetLocalizedValidationMessage("username_required", lang))
		case "Password":
			if err.Tag() == "passwordstrength" {
				messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
			} else {
				messages = append(messages, localization.GetLocalizedValidationMessage("password_required", lang))
			}
		case "Email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_required", lang))
		case "PhoneNumber":
			messages = append(messages, localization.GetLocalizedValidationMessage("phone_number_required", lang))
		}
	}
	return messages
}
