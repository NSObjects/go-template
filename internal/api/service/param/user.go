/*
 * Generated from OpenAPI3 document
 * Module: user
 */

package param




// UserGetByIDRequest 请求结构体
type UserGetByIDRequest struct {
	ID int64 `json:"id" form:"id" param:"id" validate:"required,min=1"`
}

// UserUpdateRequest 请求结构体
type UserUpdateRequest struct {
	Id int `json:"id" validate:"min=1"`
	Name string `json:"name" validate:"required,min=2,max=50"`
	Phone string `json:"phone" validate:"required,regexp=^1[3-9]\d{9}$"`
	Account string `json:"account" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=128"`
	Status int `json:"status" validate:"min=0,max=1"`
}

// UserDeleteRequest 请求结构体
type UserDeleteRequest struct {
	ID int64 `json:"id" param:"id" validate:"required,min=1"`
}

// UserListRequest 请求结构体
type UserListRequest struct {
	Page  int    `json:"page" form:"page" query:"page" validate:"min=1"`
	Count int    `json:"count" form:"count" query:"count" validate:"min=1,max=100"`
	Name string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

// UserCreateRequest 请求结构体
type UserCreateRequest struct {
	Id int `json:"id" validate:"min=1"`
	Name string `json:"name" validate:"required,min=2,max=50"`
	Phone string `json:"phone" validate:"required,regexp=^1[3-9]\d{9}$"`
	Account string `json:"account" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=128"`
	Status int `json:"status" validate:"min=0,max=1"`
}

// UserResponse 响应结构体
type UserResponse struct {
	// TODO: 根据OpenAPI文档定义响应字段
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UserParam 查询参数结构体
type UserParam struct {
	Page int `json:"page" form:"page" query:"page"`
	Count int `json:"count" form:"count" query:"count"`
	Name string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}
