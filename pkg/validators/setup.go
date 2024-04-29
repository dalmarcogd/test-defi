package validators

import (
	"github.com/go-playground/validator/v10"
)

func Setup() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("notblank", NotBlank)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("notspace", NotSpace)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("cpf", CPF)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("cnpj", CNPJ)
	if err != nil {
		return nil, err
	}

	return validate, nil
}
