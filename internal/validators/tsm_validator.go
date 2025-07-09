package validators

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

var validate = validator.New()

func ValidateTsmRegister(req *dtos.TsmRegisterRequest) []string {
	if err := validate.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make([]string, 0, len(ve))
			for _, fe := range ve {
				// Try to find a specific message key for the field and tag
				msgKey := fmt.Sprintf("validation_%s_%s", strings.ToLower(fe.Field()), strings.ToLower(fe.Tag()))
				
				// If no specific message key, use the generic one
				// The handler will then localize and format this generic message
				if !isSpecificMessageKeyDefined(msgKey) {
					msgKey = "validation_generic_field_failed"
				}
				errMsgs = append(errMsgs, msgKey)
			}
			return errMsgs
		}
		return []string{"validation_generic_error"}
	}
	return nil
}

// isSpecificMessageKeyDefined is a helper to check if a specific message key exists.
// In a real application, this would check a map of all defined message keys.
// For simplicity here, we assume if it's not a generic one, it's specific.
func isSpecificMessageKeyDefined(key string) bool {
	switch key {
	case "validation_appcode_required",
		"validation_merchantcode_required",
		"validation_terminalcode_required",
		"validation_serialnumber_required",
		"validation_mid_required":
		return true
	default:
		return false
	}
}