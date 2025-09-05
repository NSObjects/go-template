package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ./error_code_generated.md

// Order相关错误码
const (
	// ErrOrderNotFound - 404: Order not found.
	ErrOrderNotFound int = iota + 212140
	// ErrOrderAlreadyExists - 400: Order already exists.
	ErrOrderAlreadyExists
	// ErrOrderInvalidData - 400: Order invalid data.
	ErrOrderInvalidData
	// ErrOrderPermissionDenied - 403: Order permission denied.
	ErrOrderPermissionDenied
	// ErrOrderInUse - 400: Order is in use.
	ErrOrderInUse
	// ErrOrderCreateFailed - 500: Order create failed.
	ErrOrderCreateFailed
	// ErrOrderUpdateFailed - 500: Order update failed.
	ErrOrderUpdateFailed
	// ErrOrderDeleteFailed - 500: Order delete failed.
	ErrOrderDeleteFailed
)
