package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

func ValidateProductAddOnRequest(s *dtos.ProductAddOnRequest) error {
	validate := validator.New()
	return validate.Struct(s)
}
