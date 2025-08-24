package dto

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Error error
	Field string
}

var v = validator.New()

func validate(dto any, tagMessages map[string]string) *ValidationError {
	err := v.Struct(dto)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return &ValidationError{Error: err}
	}

	return getError(validationErrors, tagMessages)
}

func getError(validationErrors validator.ValidationErrors, tagMessages map[string]string) *ValidationError {
	for _, fieldError := range validationErrors {
		field := fieldError.Field()

		if msg, exists := tagMessages[fieldError.Tag()]; exists {
			field = strings.ToLower(field)

			return &ValidationError{Error: errors.New("Поле " + field + " " + msg), Field: field}
		}
	}

	return nil
}
