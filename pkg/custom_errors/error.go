package custom_errors

import (
	"fmt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type responseError struct {
	errorType     ErrorType
	originalError error
	contextInfo   errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// Error returns the message of a responseError
func (error responseError) Error() string {
	return error.originalError.Error()
}

// New creates a new responseError
func (errorType ErrorType) New(span trace.Span, recordSpan bool, msg string) error {
	err := responseError{
		errorType:     errorType,
		originalError: errors.New(msg),
	}
	if recordSpan {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

	}
	return err
}

// Newf creates a new responseError with formatted message
func (errorType ErrorType) Newf(span trace.Span, recordSpan bool, msg string, args ...interface{}) error {
	err := responseError{
		errorType:     errorType,
		originalError: fmt.Errorf(msg, args...),
	}
	if span != nil && recordSpan {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(span trace.Span, recordSpan bool, err error, msg string) error {
	return errorType.Wrapf(span, recordSpan, err, msg)
}

// Wrapf creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(span trace.Span, recordSpan bool, err error, msg string, args ...interface{}) error {
	wrapErr := responseError{
		errorType:     errorType,
		originalError: errors.Wrapf(err, msg, args...),
	}
	if span != nil && recordSpan {
		span.RecordError(wrapErr)
		span.SetStatus(codes.Error, err.Error())

	}
	return wrapErr
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if responseErr, ok := err.(responseError); ok {
		return responseError{
			errorType:     responseErr.errorType,
			originalError: responseErr.originalError,
			contextInfo:   context,
		}
	}

	return responseError{
		errorType:     InternalError,
		originalError: err,
		contextInfo:   context,
	}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if responseErr, ok := err.(responseError); ok && responseErr.contextInfo != emptyContext {
		return map[string]string{
			"field":   responseErr.contextInfo.Field,
			"message": responseErr.contextInfo.Message,
		}
	}

	return nil
}

// GetErrorType returns the error type
func GetErrorType(err error) ErrorType {
	if responseErr, ok := err.(responseError); ok {
		return responseErr.errorType
	}

	return InternalError
}

func AddRequestContextToError(reqInfo string, err error) error {
	errorType := GetErrorType(err)
	status := GetHttpStatusCode(errorType)

	if status != 420 {
		var ctxError error
		if status == 400 {
			ctxError = BadRequest.Newf(nil, false, "Request Info:%s \n, Error: %s ", reqInfo, err.Error())
		} else {
			ctxError = InternalError.Newf(nil, false, "Request Info:%s \n, Error: %s ", reqInfo, err.Error())

		}
		return ctxError
	} else {
		return err
	}
}
