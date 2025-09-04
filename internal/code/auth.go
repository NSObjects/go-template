package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ./error_code_generated.md

// Auth相关错误码
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
	// ErrAuthCreateFailed - 500: Auth create failed.
	ErrAuthCreateFailed
	// ErrAuthUpdateFailed - 500: Auth update failed.
	ErrAuthUpdateFailed
	// ErrAuthDeleteFailed - 500: Auth delete failed.
	ErrAuthDeleteFailed
)
