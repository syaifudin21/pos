package routes

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	internalmw "github.com/msyaifudin/pos/internal/middleware"
)

func WithValidation(dtoType interface{}, validatorFunc interface{}) echo.MiddlewareFunc {
	return internalmw.ValidationMiddleware(dtoType, func(data interface{}) interface{} {
		// Use reflection to call the actual validator function with the correct type
		validatorValue := reflect.ValueOf(validatorFunc)
		// Ensure data is a pointer if the validator expects a pointer
		var arg reflect.Value
		if reflect.TypeOf(dtoType).Kind() == reflect.Ptr {
			arg = reflect.ValueOf(data)
		} else {
			arg = reflect.ValueOf(data).Elem()
		}

		results := validatorValue.Call([]reflect.Value{arg})
		if len(results) > 0 && !results[0].IsNil() {
			// Handle validator.ValidationErrors specifically
			if ve, ok := results[0].Interface().(validator.ValidationErrors); ok {
				var messages []string
				for _, err := range ve {
					messages = append(messages, err.Field()+":"+err.Tag())
				}
				return messages
			} else if msgs, ok := results[0].Interface().([]string); ok {
				return msgs
			}
		}
		return nil
	})
}
