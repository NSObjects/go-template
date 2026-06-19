package apperr

import "net/http"

// HTTPStatus returns the HTTP status associated with a registered application
// error code. Unknown codes are treated as internal server errors.
func HTTPStatus(code int) int {
	if code == 0 {
		return http.StatusOK
	}

	def, ok := Lookup(code)
	if !ok {
		return http.StatusInternalServerError
	}

	switch def.Kind {
	case KindOK:
		return http.StatusOK
	case KindBadRequest, KindValidation:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// IsClientError reports whether code maps to a 4xx HTTP status.
func IsClientError(code int) bool {
	status := HTTPStatus(code)
	return status >= http.StatusBadRequest && status < http.StatusInternalServerError
}

// IsServerError reports whether code maps to a 5xx HTTP status.
func IsServerError(code int) bool {
	return HTTPStatus(code) >= http.StatusInternalServerError
}

// IsInternalError reports whether code should be handled as an internal error.
func IsInternalError(code int) bool {
	return IsServerError(code)
}
