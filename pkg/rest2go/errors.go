package rest2go

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type FieldError struct {
	Field    string `json:"field"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Value    string `json:"value"`
	Expected string `json:"expected"`
}

func NewFieldError(field, code, message string) FieldError {
	return FieldError{
		Field:   field,
		Code:    code,
		Message: message,
	}
}

func NewDetailedFieldError(field, code, message, value, expected string) FieldError {
	return FieldError{
		Field:    field,
		Code:     code,
		Message:  message,
		Value:    value,
		Expected: expected,
	}
}

type ApiErrorDto struct {
	Timestamp string       `json:"timestamp"`
	Code      string       `json:"code"`
	Message   string       `json:"message"`
	Details   []FieldError `json:"details"`
}

func NewApiErrorDto(code, message string, details ...FieldError) ApiErrorDto {
	return ApiErrorDto{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Code:      code,
		Message:   message,
		Details:   details,
	}
}

type ApiError struct {
	Status  int
	Cause   string
	Details []FieldError
}

func (e *ApiError) Error() string {
	return e.Cause
}

func NewApiError(status int, cause string, details ...FieldError) *ApiError {
	if len(details) == 0 {
		details = []FieldError{}
	}

	return &ApiError{
		Status:  status,
		Cause:   cause,
		Details: details,
	}
}

func dtoToJson(dto ApiErrorDto) []byte {
	bytes, _ := json.Marshal(dto)

	return bytes
}

func HandleError(err error, response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
	var apiErr *ApiError

	if errors.As(err, &apiErr) {
		slog.Error("Handling expected API error", "err", err.Error())
		var dto ApiErrorDto

		switch apiErr.Status {
		case 400:
			dto = NewApiErrorDto("err.invalid-request", "Invalid Request", apiErr.Details...)

		case 401:
			dto = NewApiErrorDto("auth.unauthorized", "Unauthorized")

		case 403:
			dto = NewApiErrorDto("auth.access-denied", "Access Denied")

		case 404:
			dto = NewApiErrorDto("err.not-found", "Not Found")

		case 406:
			dto = NewApiErrorDto("err.not-acceptable", "Not Acceptable")

		case 500:
			dto = NewApiErrorDto("err.internal", "Internal Server Error", apiErr.Details...)

		default:
			dto = NewApiErrorDto("err.unknown", "Unknown Server Error", apiErr.Details...)
		}

		response.WriteHeader(apiErr.Status)
		response.Write(dtoToJson(dto))
		return
	}

	slog.Error("Handling unexpected API error", "err", err.Error())
	errDto := NewApiErrorDto("err.internal", "Internal Server Error")
	response.WriteHeader(http.StatusInternalServerError)
	response.Write(dtoToJson(errDto))
}

func HandleNotFoundError(response http.ResponseWriter, request *http.Request) {
	HandleError(NewApiError(404, "not found"), response)
}
