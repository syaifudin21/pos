package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
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

func ValidateRegisterRequest(req *dtos.RegisterRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		case "passwordstrength":
			messages = append(messages, "password_strength")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateLoginRequest(req *dtos.LoginRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		case "passwordstrength":
			messages = append(messages, "password_strength")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateUpdateUserRequest(req *dtos.UpdateUserRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		case "passwordstrength":
			messages = append(messages, "password_strength")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateRegisterAdminRequest(req *dtos.RegisterAdminRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		case "passwordstrength":
			messages = append(messages, "password_strength")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateVerifyOTPRequest(req *dtos.VerifyOTPRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		case "len":
			messages = append(messages, err.Field()+"_invalid")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateUpdatePasswordRequest(req *dtos.UpdatePasswordRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "passwordstrength":
			messages = append(messages, "password_strength")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateSendOTPRequest(req *dtos.SendOTPRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateUpdateEmailRequest(req *dtos.UpdateEmailRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		case "len":
			messages = append(messages, err.Field()+"_invalid")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateForgotPasswordRequest(req *dtos.ForgotPasswordRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateResetPasswordRequest(req *dtos.ResetPasswordRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		case "len":
			messages = append(messages, err.Field()+"_invalid")
		case "passwordstrength":
			messages = append(messages, "password_strength")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}

func ValidateResendEmailRequest(req *dtos.ResendEmailRequest) []string {
	err := authValidator.Struct(req)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+"_required")
		case "email":
			messages = append(messages, "email_invalid")
		default:
			messages = append(messages, err.Field()+"_"+err.Tag())
		}
	}
	return messages
}