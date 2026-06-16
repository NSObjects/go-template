// Package apperr defines framework-free application errors.
package apperr

// Kind describes the caller-visible class of an application error without
// depending on HTTP.
type Kind string

const (
	KindOK           Kind = "ok"
	KindBadRequest   Kind = "bad_request"
	KindValidation   Kind = "validation"
	KindUnauthorized Kind = "unauthorized"
	KindForbidden    Kind = "forbidden"
	KindNotFound     Kind = "not_found"
	KindConflict     Kind = "conflict"
	KindInternal     Kind = "internal"
)

// Category identifies the operational area associated with an error.
type Category string

const (
	CategorySystem     Category = "system"
	CategoryDatabase   Category = "database"
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
)

// Data source and external system error codes.
const (
	ErrDatabase int = iota + 100101
	ErrExternalService
)

// Common request error codes.
const (
	ErrBadRequest int = iota + 100400
	ErrUnauthorized
	ErrForbidden
	ErrNotFound
	ErrInternalServer
)

// Authentication error codes used by the HTTP runtime.
const (
	ErrSignatureInvalid int = iota + 100201
)

// Definition is the registered meaning of an application error code.
type Definition struct {
	Code     int
	Kind     Kind
	Category Category
	Message  string
}

var definitions = map[int]Definition{
	ErrSuccess:          {Code: ErrSuccess, Kind: KindOK, Category: CategorySystem, Message: "OK"},
	ErrUnknown:          {Code: ErrUnknown, Kind: KindInternal, Category: CategorySystem, Message: "Internal server error"},
	ErrBind:             {Code: ErrBind, Kind: KindBadRequest, Category: CategoryValidation, Message: "Error occurred while binding the request body to the struct"},
	ErrValidation:       {Code: ErrValidation, Kind: KindValidation, Category: CategoryValidation, Message: "Validation failed"},
	ErrDatabase:         {Code: ErrDatabase, Kind: KindInternal, Category: CategoryDatabase, Message: "Database error"},
	ErrExternalService:  {Code: ErrExternalService, Kind: KindInternal, Category: CategoryExternal, Message: "External service error"},
	ErrBadRequest:       {Code: ErrBadRequest, Kind: KindBadRequest, Category: CategoryValidation, Message: "Bad request"},
	ErrUnauthorized:     {Code: ErrUnauthorized, Kind: KindUnauthorized, Category: CategoryAuth, Message: "Unauthorized"},
	ErrForbidden:        {Code: ErrForbidden, Kind: KindForbidden, Category: CategoryPermission, Message: "Forbidden"},
	ErrNotFound:         {Code: ErrNotFound, Kind: KindNotFound, Category: CategoryBusiness, Message: "Not found"},
	ErrInternalServer:   {Code: ErrInternalServer, Kind: KindInternal, Category: CategorySystem, Message: "Internal server error"},
	ErrSignatureInvalid: {Code: ErrSignatureInvalid, Kind: KindUnauthorized, Category: CategoryAuth, Message: "Signature is invalid"},
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
