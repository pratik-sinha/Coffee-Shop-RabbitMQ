package custom_errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	BadRequest = ErrorType(iota)
	Unauthorized
	Forbidden
	NotFound
	Conflict
	InternalError
	Unavailable
	UIError
)

type ErrorType uint

// GetStatusCode  the status code for the error type
func GetHttpStatusCode(errorType ErrorType) int {
	switch errorType {
	case BadRequest:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
	case InternalError:
		return http.StatusInternalServerError
	case Unavailable:
		return http.StatusServiceUnavailable
	case UIError:
		return 420
	default:
		return http.StatusInternalServerError
	}
}

func GetGrpcStatusCode(errorType ErrorType) codes.Code {
	switch errorType {
	case BadRequest:
		return codes.InvalidArgument
	case Unauthorized:
		return codes.Unauthenticated
	case Forbidden:
		return codes.PermissionDenied
	case NotFound:
		return codes.NotFound
	case InternalError:
		return codes.Internal
	case Unavailable:
		return codes.Unavailable
	case UIError:
		return 420
	default:
		return codes.Internal
	}
}
