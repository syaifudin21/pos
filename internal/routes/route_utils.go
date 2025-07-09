package routes

import (
	"reflect"

	"github.com/labstack/echo/v4"
	internalmw "github.com/msyaifudin/pos/internal/middleware"
)

func WithValidation(dtoType interface{}, validatorFunc interface{}) echo.MiddlewareFunc {
	return internalmw.ValidationMiddleware(dtoType, func(data interface{}) []string {
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
			return results[0].Interface().([]string)
		}
		return nil
	})
}
