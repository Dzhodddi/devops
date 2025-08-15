package main

import (
	"github.com/go-playground/validator/v10"
)

type ValidatorDirective struct {
	Validate *validator.Validate
}

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
