package errorx

import (
	"encoding/json"
	"errors"

	"go.uber.org/zap/zapcore"
)

// BusinessError defines the user-facing part of an error.
type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DetailError is the full error object used internally.
type DetailError struct {
	BusinessError
	HTTPStatus int
	LogLevel   zapcore.Level
	Err        error // The original underlying error
}

func (e *DetailError) Error() string {
	return e.Err.Error()
}

// Unwrap for errors.Is and errors.As
func (e *DetailError) Unwrap() error {
	return e.Err
}

// ToJSON converts the BusinessError to a JSON string.
func (be *BusinessError) ToJSON() string {
	s, _ := json.Marshal(be)
	return string(s)
}

// New creates a new DetailError.
// It looks up the code in the predefined map to get the message and HTTP status.
func New(code int, underlyingErr error) *DetailError {
	be, ok := codes[code]
	if !ok {
		// Fallback for undefined codes
		be = codes[CodeInternalServerError]
	}

	// If no specific underlying error is provided, use the message from the map.
	err := underlyingErr
	if err == nil {
		err = errors.New(be.Message)
	}

	return &DetailError{
		BusinessError: BusinessError{
			Code:    be.Code,
			Message: be.Message,
		},
		HTTPStatus: be.HTTPStatus,
		LogLevel:   be.LogLevel,
		Err:        err,
	}
}

// WithMessage allows overriding the default message for a given error code.
func (e *DetailError) WithMessage(message string) *DetailError {
	e.BusinessError.Message = message
	return e
}
