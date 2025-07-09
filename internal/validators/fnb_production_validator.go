package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

func ValidateFNBProductionRequest(s *dtos.FNBProductionRequest) error {
	validate := validator.New()
	return validate.Struct(s)
}
