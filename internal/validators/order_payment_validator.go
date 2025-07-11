package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

func ValidateCreateOrderPayment(req *dtos.CreateOrderPaymentRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}
