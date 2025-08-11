package errorx

import (
	"net/http"

	"go.uber.org/zap/zapcore"
)

// Predefined error codes.
const (
	// Common Errors (10xxx)
	CodeSuccess             = 0
	CodeInternalServerError = 10001
	CodeInvalidParams       = 10002
	CodeUnauthorized        = 10003
	CodeNotFound            = 10004

	// User Service Errors (20xxx)
	CodeUserAlreadyExists  = 20001
	CodeUserNotFound       = 20002
	CodeInvalidCredentials = 20003

	// Article Service Errors (30xxx)
	CodeArticleNotFound = 30001
)

var codes = map[int]struct {
	Code       int
	Message    string
	HTTPStatus int
	LogLevel   zapcore.Level
}{
	CodeSuccess:             {CodeSuccess, "Success", http.StatusOK, zapcore.InfoLevel},
	CodeInternalServerError: {CodeInternalServerError, "Internal Server Error", http.StatusInternalServerError, zapcore.ErrorLevel},
	CodeInvalidParams:       {CodeInvalidParams, "Invalid Parameters", http.StatusBadRequest, zapcore.WarnLevel},
	CodeUnauthorized:        {CodeUnauthorized, "Unauthorized", http.StatusUnauthorized, zapcore.WarnLevel},
	CodeNotFound:            {CodeNotFound, "Resource Not Found", http.StatusNotFound, zapcore.WarnLevel},
	CodeUserAlreadyExists:   {CodeUserAlreadyExists, "User already exists", http.StatusBadRequest, zapcore.WarnLevel},
	CodeUserNotFound:        {CodeUserNotFound, "User not found", http.StatusNotFound, zapcore.WarnLevel},
	CodeInvalidCredentials:  {CodeInvalidCredentials, "Invalid username or password", http.StatusUnauthorized, zapcore.WarnLevel},
	CodeArticleNotFound:     {CodeArticleNotFound, "Article not found", http.StatusNotFound, zapcore.WarnLevel},
}
