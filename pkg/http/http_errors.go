package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gallery-service/pkg/constants"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
)

const (
	ErrBadRequest          = "Bad request"
	ErrNotFound            = "Not Found"
	ErrUnauthorized        = "Unauthorized"
	ErrRequestTimeout      = "Request Timeout"
	ErrInvalidProductName  = "Invalid cluster name"
	ErrInvalidSBCode       = "Invalid SB-Code"
	ErrInvalidField        = "Invalid field"
	ErrInvalidID           = "Invalid id provided"
	ErrInternalServerError = "Internal Server Error"
)

var (
	BadRequest          = errors.New("Bad request")
	WrongCredentials    = errors.New("Wrong Credentials")
	NotFound            = errors.New("Not Found")
	Unauthorized        = errors.New("Unauthorized")
	Forbidden           = errors.New("Forbidden")
	TooManyRequest      = errors.New("Too Many Request")
	InternalServerError = errors.New("Internal Server Error")
)

// RestErr Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() string
	ErrBody() RestError
}

// RestError Rest error struct
type RestError struct {
	ErrStatus  int       `json:"status_code,omitempty"`
	ErrError   string    `json:"error,omitempty"`
	ErrMessage string    `json:"message,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
}

// ErrBody Error body
func (e RestError) ErrBody() RestError {
	return e
}

// Error  Error() interface method
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrMessage)
}

// Status Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// Causes RestError Causes
func (e RestError) Causes() string {
	return e.ErrMessage
}

// NewRestError New Rest Error
func NewRestError(status int, err string, causes string, debug bool) RestErr {
	restError := RestError{
		ErrStatus: status,
		ErrError:  err,
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return restError
}

// NewRestErrorWithMessage New Rest Error With MsgMessage
func NewRestErrorWithMessage(status int, err string, causes string) RestErr {
	return RestError{
		ErrStatus:  status,
		ErrError:   err,
		ErrMessage: causes,
		Timestamp:  time.Now().UTC(),
	}
}

// NewRestErrorFromBytes New Rest Error From Bytes
func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr RestError
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

// NewBadRequestError New Bad Request Error
func NewBadRequestError(ctx *fiber.Ctx, causes string, debug bool) error {
	restError := RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  BadRequest.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	restErr, err := json.Marshal(restError)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusBadRequest, string(restErr))
}

// NewNotFoundError New Not Found Error
func NewNotFoundError(ctx *fiber.Ctx, causes string, debug bool) error {
	restError := RestError{
		ErrStatus: http.StatusNotFound,
		ErrError:  NotFound.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	restErr, err := json.Marshal(restError)

	if err != nil {
		return ctx.JSON(http.StatusNotFound, err.Error())
	}
	return ctx.JSON(http.StatusNotFound, string(restErr))
}

// NewUnauthorizedError New Unauthorized Error
func NewUnauthorizedError(ctx *fiber.Ctx, causes string, debug bool) error {

	restError := RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  Unauthorized.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	restErr, err := json.Marshal(restError)

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, err.Error())
	}
	return ctx.JSON(http.StatusUnauthorized, string(restErr))
}

// NewForbiddenError New Forbidden Error
func NewForbiddenError(ctx *fiber.Ctx, causes string, debug bool) error {

	restError := RestError{
		ErrStatus: http.StatusForbidden,
		ErrError:  Forbidden.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	restErr, err := json.Marshal(restError)

	if err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	return ctx.JSON(http.StatusForbidden, string(restErr))
}

// NewTooManyRequestError New TooManyRequest Error
func NewTooManyRequestError(ctx *fiber.Ctx, causes string, debug bool) error {
	restError := RestError{
		ErrStatus: http.StatusTooManyRequests,
		ErrError:  TooManyRequest.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	restErr, err := json.Marshal(restError)

	if err != nil {
		return ctx.JSON(http.StatusTooManyRequests, err.Error())
	}
	return ctx.JSON(http.StatusTooManyRequests, string(restErr))
}

// NewInternalServerError New Internal Server Error
func NewInternalServerError(ctx *fiber.Ctx, causes string, debug bool) error {

	restError := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	restErr, err := json.Marshal(restError)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusInternalServerError, string(restErr))
}

// ParseErrors Parser of error string messages returns RestError
func ParseErrors(err error, debug bool) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, ErrNotFound, err.Error(), debug)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, ErrRequestTimeout, err.Error(), debug)
	case errors.Is(err, Unauthorized):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case errors.Is(err, WrongCredentials):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)

	case strings.Contains(strings.ToLower(err.Error()), constants.SQLState):
		return parseSqlErrors(err, debug)
	case strings.Contains(strings.ToLower(err.Error()), "field validation"):
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return NewRestError(http.StatusBadRequest, ErrBadRequest, validationErrors.Error(), debug)
		}
		return parseValidatorError(err, debug)

	case strings.Contains(strings.ToLower(err.Error()), "required header"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), constants.Base64):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), constants.Unmarshal):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), constants.Uuid):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), constants.Cookie):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), constants.Token):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), constants.Bcrypt):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)

	case strings.Contains(strings.ToLower(err.Error()), "no documents in result"):
		return NewRestError(http.StatusNotFound, ErrNotFound, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "not found"):
		return NewRestError(http.StatusNotFound, ErrNotFound, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "already in use"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "not a valid objectid"):
		return NewRestError(http.StatusBadRequest, ErrInvalidID, err.Error(), debug)

	default:
		if err, ok := err.(*RestError); ok {
			return err
		}
		return NewRestError(http.StatusInternalServerError, ErrInternalServerError, errors.Cause(err).Error(), debug)
	}

}

func parseSqlErrors(err error, debug bool) RestErr {
	return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
}

func parseValidatorError(err error, debug bool) RestErr {
	if strings.Contains(err.Error(), "ProductName") {
		return NewRestError(http.StatusBadRequest, ErrInvalidProductName, err.Error(), debug)
	}

	if strings.Contains(err.Error(), "SBCode") {
		return NewRestError(http.StatusBadRequest, ErrInvalidSBCode, err.Error(), debug)
	}

	return NewRestError(http.StatusBadRequest, ErrInvalidField, err.Error(), debug)
}

// ErrorResponse Error response
func ErrorResponse(err error, debug bool) (int, interface{}) {
	return ParseErrors(err, debug).Status(), ParseErrors(err, debug)
}

// ErrorCtxResponse Error response object and status code
func ErrorCtxResponse(ctx *fiber.Ctx, err error, debug bool) error {
	restErr := ParseErrors(err, debug)

	return ctx.Status(restErr.Status()).JSON(restErr)
}
