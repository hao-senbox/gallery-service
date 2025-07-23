package validator

import (
	"fmt"
	"gallery-service/config"
	"gallery-service/pkg/zap"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"strings"
)

// Wrapper implementation for Fiber
type Wrapper struct {
	log       zap.Logger
	cfg       *config.Config
	validator *validator.Validate
}

// New creates a new validator instance
func NewValidator(log zap.Logger, cfg *config.Config) *Wrapper {
	return &Wrapper{
		log:       log,
		cfg:       cfg,
		validator: validator.New(),
	}
}

type ErrorResponse struct {
	Error       bool        `json:"error"`
	FailedField string      `json:"failed_field,omitempty"`
	Tag         string      `json:"tag,omitempty"`
	Value       interface{} `json:"value,omitempty"`
}

// Validate performs validation on the provided data
func (w *Wrapper) Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse
	err := w.validator.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ErrorResponse{
				Error:       true,
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Value(),
			})
		}
	}

	return validationErrors
}

// DataValidation returns a Fiber middlewares for validating request data
func (w *Wrapper) DataValidation(data interface{}) error {
	// Validate the parsed data
	validationErrors := w.Validate(data)
	if len(validationErrors) > 0 {
		errMsgs := make([]string, 0)
		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}
		w.log.Errorf("(validate) err: {%v}", errMsgs)

		return errors.New(fmt.Sprintf("invalid field validation: %s", strings.Join(errMsgs, " and ")))
	}

	return nil
}
