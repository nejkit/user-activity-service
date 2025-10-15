package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"user-activity-service/services"
)

var (
	invalidValidationError *validator.InvalidValidationError
)

func validateRequest[T any](req *T) error {
	valCtx := validator.New()

	err := valCtx.Struct(req)

	return processValidationResult(err)
}

func processValidationResult(err error) error {
	if err == nil {
		return nil
	}

	if errors.As(err, &invalidValidationError) {
		return services.ErrorInternalError
	}

	for _, validatorErr := range err.(validator.ValidationErrors) {
		return errors.New(validatorErr.Error())
	}

	return nil
}
