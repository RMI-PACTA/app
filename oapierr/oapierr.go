package oapierr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logLevel int

// These log levels mirror the levels available in zapcore.Level, since that's
// the structured logger we've chosen.
const (
	unsetLevel logLevel = iota
	debugLevel
	infoLevel
	warnLevel
	errorLevel
	panicLevel
)

var (
	// Generally, things that could ostensibly be triggered by a client are
	// warnings, to cut down on noise/pager fatigue. Errors are almost always
	// programmer error or backend system failure.
	defaultLevelForCode = map[int]logLevel{
		http.StatusBadRequest:          warnLevel,
		http.StatusNotFound:            warnLevel,
		http.StatusConflict:            warnLevel,
		http.StatusForbidden:           warnLevel,
		http.StatusTooManyRequests:     warnLevel,
		http.StatusNotImplemented:      warnLevel,
		http.StatusInternalServerError: errorLevel,
		http.StatusUnauthorized:        warnLevel,
	}

	defaultMessageForCode = map[int]string{
		http.StatusBadRequest:          "bad request",
		http.StatusNotFound:            "not found",
		http.StatusConflict:            "conflict",
		http.StatusForbidden:           "forbidden",
		http.StatusTooManyRequests:     "too many requests",
		http.StatusNotImplemented:      "not implemented",
		http.StatusInternalServerError: "internal server error",
		http.StatusUnauthorized:        "unauthorized",
	}
)

// ErrorID represents a type of error specific to the domain of the caller of
// this package, like admin_only or too_many_muffins.
type ErrorID string

type Error struct {
	// These must be set for all errors
	statusCode int
	msg        string

	// A default log level will be chosen based on
	// code if none is provided.
	level logLevel
	// A default client message will be chosen based
	// on the code if none is provided.
	clientMsg string

	// These are additional metadata that do not need
	// to be set.
	fields  []zap.Field
	errorID ErrorID
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	err := e.err()
	if err == nil {
		// We only write the code and message for now, the actual logger should log
		// the fields.
		return fmt.Sprintf("[%q] %s", e.statusCode, e.msg)
	}
	return fmt.Sprintf("[%q] %s: %v", e.statusCode, e.msg, err)
}

func (e *Error) err() error {
	for _, f := range e.fields {
		if f.Key != "error" || f.Type != zapcore.ErrorType {
			continue
		}
		errVal, ok := f.Interface.(error)
		if !ok {
			continue
		}
		return errVal
	}
	return nil
}

func (e *Error) logLevel() logLevel {
	// Return a level if one was explicitly set
	if e.level != unsetLevel {
		return e.level
	}

	// Otherwise, return whatever the default level is for that code.
	level, ok := defaultLevelForCode[e.statusCode]
	if !ok {
		// We didn't have a default for it, which seems pretty exceptional.
		return errorLevel
	}

	return level
}

func (e *Error) ErrorID() ErrorID {
	return e.errorID
}

func (e *Error) ClientMessage() (string, bool) {
	// Return the client message if one was explicitly set.
	if e.clientMsg != "" {
		return e.clientMsg, true
	}

	// Otherwise, return whatever the default is for that code.
	return defaultMessageForCode[e.statusCode], false
}

// WithMessage adds an error intended for clients to see, and returns the error
// for chaining purposes.
func (e *Error) WithMessage(msg string) *Error {
	e.clientMsg = msg
	return e
}

// WithErrorID adds an error intended for client apps to use, and returns
// the error for chaining purposes. It'll appear in the GraphQL response
// "extensions" field, see:
// https://spec.graphql.org/October2021/#sec-Response-Format
func (e *Error) WithErrorID(errID ErrorID) *Error {
	e.errorID = errID
	return e
}

// AtDebug overrides the default level for the error and logs at DEBUG level.
func (e *Error) AtDebug() *Error {
	e.level = debugLevel
	return e
}

// AtInfo overrides the default level for the error and logs at INFO level.
func (e *Error) AtInfo() *Error {
	e.level = infoLevel
	return e
}

// AtWarn overrides the default level for the error and logs at WARN level.
func (e *Error) AtWarn() *Error {
	e.level = warnLevel
	return e
}

// AtError overrides the default level for the error and logs at ERROR level.
func (e *Error) AtError() *Error {
	e.level = errorLevel
	return e
}

// AtPanic overrides the default level for the error and logs at PANIC level.
// Note: This doesn't actually cause a panic when using a production zap logger,
// since we use DPanic.
func (e *Error) AtPanic() *Error {
	e.level = panicLevel
	return e
}

// New returns an initialize error with the given code. The message and fields
// are used for logging, and won't be visible to clients. For setting client-
// visible response parameters, see WithErrorID and WithMessage
func New(code int, msg string, fields ...zap.Field) *Error {
	return &Error{
		statusCode: code,
		msg:        msg,
		fields:     fields,
	}
}

func BadRequest(msg string, fields ...zap.Field) *Error {
	return New(http.StatusBadRequest, msg, fields...)
}
func NotFound(msg string, fields ...zap.Field) *Error {
	return New(http.StatusNotFound, msg, fields...)
}
func Conflict(msg string, fields ...zap.Field) *Error {
	return New(http.StatusConflict, msg, fields...)
}
func Forbidden(msg string, fields ...zap.Field) *Error {
	return New(http.StatusForbidden, msg, fields...)
}
func TooManyRequests(msg string, fields ...zap.Field) *Error {
	return New(http.StatusTooManyRequests, msg, fields...)
}
func NotImplemented(msg string, fields ...zap.Field) *Error {
	return New(http.StatusNotImplemented, msg, fields...)
}
func Internal(msg string, fields ...zap.Field) *Error {
	return New(http.StatusInternalServerError, msg, fields...)
}
func Unauthorized(msg string, fields ...zap.Field) *Error {
	return New(http.StatusUnauthorized, msg, fields...)
}

type ResponseConverter[T any] func(*Error) T

func ErrorHandlerFunc[T any](logger *zap.Logger, fn ResponseConverter[T]) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err == nil {
			return
		}

		e := &Error{}
		if errors.As(err, &e) {
			logError(logger, e)
		} else {
			logger.Error(
				"received error that was not of type *oapierr.Error",
				zap.String("type", fmt.Sprintf("%T", err)),
				zap.Error(err),
			)
			e = Internal(err.Error(), zap.Error(err))
		}

		response := fn(e)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.statusCode)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error(
				"failed to encode error response body as JSON",
				zap.String("type", fmt.Sprintf("%T", response)),
				zap.Error(err),
			)
		}
	}
}

func logError(logger *zap.Logger, err *Error) {
	var logFn func(msg string, fields ...zap.Field)
	switch err.logLevel() {
	case debugLevel:
		logFn = logger.Debug
	case infoLevel:
		logFn = logger.Info
	case warnLevel:
		logFn = logger.Warn
	case errorLevel:
		logFn = logger.Error
	case panicLevel:
		logFn = logger.DPanic
	default:
		// If something went wrong finding the log level, log at errorLevel.
		logFn = logger.Error
	}

	logFn(err.msg, err.fields...)
}
