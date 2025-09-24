package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ./error_code_generated.md


// User相关错误码
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 200000
	// ErrUserAlreadyExists - 400: User already exists.
	ErrUserAlreadyExists
	// ErrUserInvalidData - 400: User invalid data.
	ErrUserInvalidData
	// ErrUserPermissionDenied - 403: User permission denied.
	ErrUserPermissionDenied
	// ErrUserInUse - 400: User is in use.
	ErrUserInUse
	// ErrUserCreateFailed - 500: User create failed.
	ErrUserCreateFailed
	// ErrUser - 500: User update failed.
	ErrUserUpdateFailed
	// ErrUserDeleteFailed - 500: User delete failed.
	ErrUserDeleteFailed
)
