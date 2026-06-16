package apperr

import "fmt"

// Info is the normalized application error information used by logs and
// delivery adapters.
type Info struct {
	Kind     Kind
	Category Category
	Code     int
	Message  string
	Detail   string
}

// NewInfo normalizes err into safe public information plus internal detail.
func NewInfo(err error) Info {
	if err == nil {
		return Info{}
	}

	info := Info{
		Kind:     KindInternal,
		Category: CategorySystem,
		Code:     ErrUnknown,
		Message:  definitionFor(ErrUnknown).Message,
		Detail:   fmt.Sprintf("%+v", err),
	}

	appErr, ok := Parse(err)
	if !ok {
		return info
	}

	def, registered := Lookup(appErr.Code())
	if !registered {
		info.Detail = appErr.Detail()
		return info
	}

	info.Kind = def.Kind
	info.Category = def.Category
	info.Code = def.Code
	info.Message = appErr.Message()
	info.Detail = appErr.Detail()
	return info
}

// IsInternal reports whether info represents an internal error.
func (i Info) IsInternal() bool {
	return i.Kind == KindInternal
}

// IsBusiness reports whether info represents a non-internal error.
func (i Info) IsBusiness() bool {
	return i.Kind != "" && i.Kind != KindInternal
}
