package validation

import (
	"github.com/go-playground/validator/v10"
)

func New(options []validator.Option) *validator.Validate {
	validate := validator.New()

	for _, option := range options {
		option(validate)
	}

	return validate
}
