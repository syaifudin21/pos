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
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		case "passwordstrength":
			messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
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
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		case "passwordstrength":
			messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
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
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		case "passwordstrength":
			messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
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
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		case "passwordstrength":
			messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
		}
	}
	return messages
}

func ValidateVerifyOTPRequest(req *dtos.VerifyOTPRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		case "len":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_invalid", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
		}
	}
	return messages
}

func ValidateUpdatePasswordRequest(req *dtos.UpdatePasswordRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "passwordstrength":
			messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
		}
	}
	return messages
}

func ValidateSendOTPRequest(req *dtos.SendOTPRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
		}
	}
	return messages
}

func ValidateUpdateEmailRequest(req *dtos.UpdateEmailRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		case "len":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_invalid", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
		}
	}
	return messages
}

func ValidateForgotPasswordRequest(req *dtos.ForgotPasswordRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
		}
	}
	return messages
}

func ValidateResetPasswordRequest(req *dtos.ResetPasswordRequest, lang string) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_required", lang))
		case "email":
			messages = append(messages, localization.GetLocalizedValidationMessage("email_invalid", lang))
		case "len":
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_invalid", lang))
		case "passwordstrength":
			messages = append(messages, localization.GetLocalizedValidationMessage("password_strength", lang))
		default:
			messages = append(messages, localization.GetLocalizedValidationMessage(err.Field()+"_"+err.Tag(), lang))
		}
	}
	return messages
}
