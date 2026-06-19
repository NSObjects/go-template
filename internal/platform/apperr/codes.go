// Package apperr defines framework-free application errors.
package apperr

// Kind describes the caller-visible class of an application error without
// depending on HTTP.
type Kind string

// Application error kinds.
const (
	KindOK               Kind = "ok"
	KindBadRequest       Kind = "bad_request"
	KindValidation       Kind = "validation"
	KindUnauthorized     Kind = "unauthorized"
	KindForbidden        Kind = "forbidden"
	KindNotFound         Kind = "not_found"
	KindMethodNotAllowed Kind = "method_not_allowed"
	KindConflict         Kind = "conflict"
	KindInternal         Kind = "internal"
)

// Category identifies the operational area associated with an error.
type Category string

// Application error categories.
const (
	CategorySystem     Category = "system"
	CategoryDatabase   Category = "database"
	CategoryRedis      Category = "redis"
	CategoryKafka      Category = "kafka"
	CategoryExternal   Category = "external"
	CategoryAuth       Category = "auth"
	CategoryPermission Category = "permission"
	CategoryValidation Category = "validation"
	CategoryBusiness   Category = "business"
)

// Common application error codes.
const (
	ErrSuccess int = iota + 100001
	ErrUnknown
	ErrBind
	ErrValidation
	ErrTokenInvalid
)

// Data source and external system error codes.
const (
	ErrDatabase int = iota + 100101
	ErrRedis
	ErrKafka
	ErrExternalService
)

// Common request error codes.
const (
	ErrBadRequest       = 100400
	ErrUnauthorized     = 100401
	ErrForbidden        = 100403
	ErrNotFound         = 100404
	ErrMethodNotAllowed = 100405
	ErrConflict         = 100409
	ErrInternalServer   = 100500
)

// Authentication and authorization error codes.
const (
	ErrEncrypt int = iota + 100201
	ErrSignatureInvalid
	ErrExpired
	ErrInvalidAuthHeader
	ErrMissingHeader
	ErrPasswordIncorrect
	ErrPermissionDenied
	ErrAccountLocked
	ErrAccountDisabled
	ErrTooManyAttempts
)

// Encoding and decoding error codes.
const (
	ErrEncodingFailed int = iota + 100301
	ErrDecodingFailed
	ErrInvalidJSON
	ErrEncodingJSON
	ErrDecodingJSON
	ErrInvalidYaml
	ErrEncodingYaml
	ErrDecodingYaml
)

// Definition is the registered meaning of an application error code.
type Definition struct {
	Code     int
	Kind     Kind
	Category Category
	Message  string
}

var definitions = map[int]Definition{
	ErrSuccess:         {Code: ErrSuccess, Kind: KindOK, Category: CategorySystem, Message: "OK"},
	ErrUnknown:         {Code: ErrUnknown, Kind: KindInternal, Category: CategorySystem, Message: "Internal server error"},
	ErrBind:            {Code: ErrBind, Kind: KindBadRequest, Category: CategoryValidation, Message: "Error occurred while binding the request body to the struct"},
	ErrValidation:      {Code: ErrValidation, Kind: KindValidation, Category: CategoryValidation, Message: "Validation failed"},
	ErrTokenInvalid:    {Code: ErrTokenInvalid, Kind: KindUnauthorized, Category: CategoryAuth, Message: "Token invalid"},
	ErrDatabase:        {Code: ErrDatabase, Kind: KindInternal, Category: CategoryDatabase, Message: "Database error"},
	ErrRedis:           {Code: ErrRedis, Kind: KindInternal, Category: CategoryRedis, Message: "Redis error"},
	ErrKafka:           {Code: ErrKafka, Kind: KindInternal, Category: CategoryKafka, Message: "Kafka error"},
	ErrExternalService: {Code: ErrExternalService, Kind: KindInternal, Category: CategoryExternal, Message: "External service error"},
	ErrBadRequest:      {Code: ErrBadRequest, Kind: KindBadRequest, Category: CategoryValidation, Message: "Bad request"},
	ErrUnauthorized:    {Code: ErrUnauthorized, Kind: KindUnauthorized, Category: CategoryAuth, Message: "Unauthorized"},
	ErrForbidden:       {Code: ErrForbidden, Kind: KindForbidden, Category: CategoryPermission, Message: "Forbidden"},
	ErrNotFound:        {Code: ErrNotFound, Kind: KindNotFound, Category: CategoryBusiness, Message: "Not found"},
	ErrMethodNotAllowed: {
		Code:     ErrMethodNotAllowed,
		Kind:     KindMethodNotAllowed,
		Category: CategoryValidation,
		Message:  "Method not allowed",
	},
	ErrConflict:       {Code: ErrConflict, Kind: KindConflict, Category: CategoryBusiness, Message: "Conflict"},
	ErrInternalServer: {Code: ErrInternalServer, Kind: KindInternal, Category: CategorySystem, Message: "Internal server error"},
	ErrEncrypt: {
		Code:     ErrEncrypt,
		Kind:     KindUnauthorized,
		Category: CategoryAuth,
		Message:  "Error occurred while encrypting the user password",
	},
	ErrSignatureInvalid: {
		Code:     ErrSignatureInvalid,
		Kind:     KindUnauthorized,
		Category: CategoryAuth,
		Message:  "Signature is invalid",
	},
	ErrExpired: {
		Code:     ErrExpired,
		Kind:     KindUnauthorized,
		Category: CategoryAuth,
		Message:  "Token expired",
	},
	ErrInvalidAuthHeader: {
		Code:     ErrInvalidAuthHeader,
		Kind:     KindUnauthorized,
		Category: CategoryAuth,
		Message:  "Invalid authorization header",
	},
	ErrMissingHeader: {
		Code:     ErrMissingHeader,
		Kind:     KindUnauthorized,
		Category: CategoryAuth,
		Message:  "The `Authorization` header was empty",
	},
	ErrPasswordIncorrect: {
		Code:     ErrPasswordIncorrect,
		Kind:     KindUnauthorized,
		Category: CategoryAuth,
		Message:  "Password was incorrect",
	},
	ErrPermissionDenied: {
		Code:     ErrPermissionDenied,
		Kind:     KindForbidden,
		Category: CategoryPermission,
		Message:  "Permission denied",
	},
	ErrAccountLocked: {
		Code:     ErrAccountLocked,
		Kind:     KindForbidden,
		Category: CategoryPermission,
		Message:  "Account is locked",
	},
	ErrAccountDisabled: {
		Code:     ErrAccountDisabled,
		Kind:     KindForbidden,
		Category: CategoryPermission,
		Message:  "Account is disabled",
	},
	ErrTooManyAttempts: {
		Code:     ErrTooManyAttempts,
		Kind:     KindForbidden,
		Category: CategoryPermission,
		Message:  "Too many login attempts",
	},
	ErrEncodingFailed: {
		Code:     ErrEncodingFailed,
		Kind:     KindInternal,
		Category: CategorySystem,
		Message:  "Encoding failed due to an error with the data",
	},
	ErrDecodingFailed: {
		Code:     ErrDecodingFailed,
		Kind:     KindInternal,
		Category: CategorySystem,
		Message:  "Decoding failed due to an error with the data",
	},
	ErrInvalidJSON: {
		Code:     ErrInvalidJSON,
		Kind:     KindInternal,
		Category: CategorySystem,
		Message:  "Data is not valid JSON",
	},
	ErrEncodingJSON: {
		Code:     ErrEncodingJSON,
		Kind:     KindInternal,
		Category: CategorySystem,
		Message:  "JSON data could not be encoded",
	},
	ErrDecodingJSON: {
		Code:     ErrDecodingJSON,
		Kind:     KindInternal,
		Category: CategorySystem,
		Message:  "JSON data could not be decoded",
	},
	ErrInvalidYaml: {
		Code:     ErrInvalidYaml,
		Kind:     KindInternal,
		Category: CategorySystem,
		Message:  "Data is not valid Yaml",
	},
	ErrEncodingYaml: {
		Code:     ErrEncodingYaml,
		Kind:     KindInternal,
		Category: CategorySystem,
		Message:  "Yaml data could not be encoded",
	},
	ErrDecodingYaml: {
		Code:     ErrDecodingYaml,
		Kind:     KindInternal,
		Category: CategorySystem,
		Message:  "Yaml data could not be decoded",
	},
}

// Lookup returns the registered definition for code.
func Lookup(code int) (Definition, bool) {
	def, ok := definitions[code]
	return def, ok
}

func definitionFor(code int) Definition {
	if def, ok := Lookup(code); ok {
		return def
	}
	return definitions[ErrUnknown]
}
