package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ./error_code_generated.md

// Product相关错误码
const (
	// ErrProductNotFound - 404: Product not found.
	ErrProductNotFound int = iota + 213160
	// ErrProductAlreadyExists - 400: Product already exists.
	ErrProductAlreadyExists
	// ErrProductInvalidData - 400: Product invalid data.
	ErrProductInvalidData
	// ErrProductPermissionDenied - 403: Product permission denied.
	ErrProductPermissionDenied
	// ErrProductInUse - 400: Product is in use.
	ErrProductInUse
	// ErrProductCreateFailed - 500: Product create failed.
	ErrProductCreateFailed
	// ErrProductUpdateFailed - 500: Product update failed.
	ErrProductUpdateFailed
	// ErrProductDeleteFailed - 500: Product delete failed.
	ErrProductDeleteFailed
)
