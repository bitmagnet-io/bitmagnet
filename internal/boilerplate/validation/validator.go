package validation

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Options []Option `group:"validator_options"`
}

type Option struct {
	Apply func(*validator.Validate) error
}

type Result struct {
	fx.Out
	Validate *validator.Validate
}

func New(p Params) (r Result, err error) {
	validate := validator.New()
	for _, option := range p.Options {
		if optionErr := option.Apply(validate); optionErr != nil {
			err = optionErr
			return
		}
	}
	return Result{
		Validate: validate,
	}, nil
}

func Options(options ...Option) Option {
	return Option{
		Apply: func(v *validator.Validate) error {
			for _, option := range options {
				if err := option.Apply(v); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
