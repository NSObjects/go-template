package code

// Auth相关错误码
//
//go:generate codegen -type=int
const (
	// ErrAuthNotFound - 404: Auth not found.
	ErrAuthNotFound int = iota + 198040
	// ErrAuthAlreadyExists - 400: Auth already exists.
	ErrAuthAlreadyExists
	// ErrAuthInvalidData - 400: Auth invalid data.
	ErrAuthInvalidData
	// ErrAuthPermissionDenied - 403: Auth permission denied.
	ErrAuthPermissionDenied
	// ErrAuthInUse - 400: Auth is in use.
	ErrAuthInUse
)
