package http

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	OK        = "OK"
	Created   = "Created"
	Accepted  = "Accepted"
	NoContent = "No Content"
)

var (
	SuccessOK        = "Success OK"
	SuccessCreated   = "Resource Created"
	SuccessAccepted  = "Request Accepted"
	SuccessNoContent = "No Content"
)

// RestSuccess Rest success interface
type RestSuccess interface {
	Status() int
	Message() string
	Data() interface{}
	SuccessBody() RestSuccess
}

// RestSuccessStruct Rest success struct
type RestSuccessStruct struct {
	StatusCode int         `json:"status_code,omitempty"`
	MsgMessage string      `json:"message,omitempty"`
	MsgData    interface{} `json:"data,omitempty"`
	Timestamp  time.Time   `json:"timestamp,omitempty"`
}

// SuccessBody Success body
func (s RestSuccessStruct) SuccessBody() RestSuccess {
	return s
}

// Message Success message
func (s RestSuccessStruct) Message() string {
	return s.MsgMessage
}

// Status Success status code
func (s RestSuccessStruct) Status() int {
	return s.StatusCode
}

// Data Success response data
func (s RestSuccessStruct) Data() interface{} {
	return s.Data
}

// NewRestSuccess New Rest Success
func NewRestSuccess(status int, message string, data interface{}) RestSuccess {
	return RestSuccessStruct{
		StatusCode: status,
		MsgMessage: message,
		MsgData:    data,
		Timestamp:  time.Now().UTC(),
	}
}

// NewOKResponse New OK Success Response
func NewOKResponse(ctx *fiber.Ctx, data interface{}) error {
	return SuccessCtxResponse(ctx, http.StatusOK, SuccessOK, data)
}

// NewCreatedResponse New Created Success Response
func NewCreatedResponse(ctx *fiber.Ctx, data interface{}) error {
	return SuccessCtxResponse(ctx, http.StatusCreated, SuccessCreated, data)
}

// NewAcceptedResponse New Accepted Success Response
func NewAcceptedResponse(ctx *fiber.Ctx, data interface{}) error {
	return SuccessCtxResponse(ctx, http.StatusAccepted, SuccessAccepted, data)
}

// NewNoContentResponse New No Content Success Response
func NewNoContentResponse(ctx *fiber.Ctx) error {
	return SuccessCtxResponse(ctx, http.StatusNoContent, SuccessNoContent, nil)
}

// SuccessResponse Success response
func SuccessResponse(status int, message string, data interface{}) (int, interface{}) {
	return status, RestSuccessStruct{
		StatusCode: status,
		MsgMessage: message,
		MsgData:    data,
		Timestamp:  time.Now().UTC(),
	}
}

// SuccessCtxResponse Success response object and status code
func SuccessCtxResponse(ctx *fiber.Ctx, status int, message string, data interface{}) error {
	successResp := RestSuccessStruct{
		StatusCode: status,
		MsgMessage: message,
		MsgData:    data,
		Timestamp:  time.Now().UTC(),
	}

	return ctx.Status(successResp.Status()).JSON(successResp)
}
