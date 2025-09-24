/*
 * Generated from OpenAPI3 document
 * Module: User
 */

package param


import (
	"time"
)




// UserListUsersRequest
// 获取用户列表

// Page 页码

// Size 每页数量


type UserListUsersRequest struct {

	Page int `param:"page" query:"page" form:"page" json:"page" xml:"page" validate:"min=1"`

	Size int `param:"size" query:"size" form:"size" json:"size" xml:"size" validate:"min=1,max=100"`


}



// UserCreateRequest
// 创建用户



// Username 用户名

// Email 邮箱

// Age 年龄


type UserCreateRequest struct {



	Username string `param:"username" query:"username" form:"username" json:"username" xml:"username" validate:"required,min=3,max=20"`

	Email string `param:"email" query:"email" form:"email" json:"email" xml:"email" validate:"required,email"`

	Age int `param:"age" query:"age" form:"age" json:"age" xml:"age" validate:"min=0,max=150"`


}





// UserUpdateRequest
// 更新用户



// Age 年龄

// Username 用户名

// Email 邮箱


type UserUpdateRequest struct {



	Age int `param:"age" query:"age" form:"age" json:"age" xml:"age" validate:"min=0,max=150"`

	Username string `param:"username" query:"username" form:"username" json:"username" xml:"username" validate:"min=3,max=20"`

	Email string `param:"email" query:"email" form:"email" json:"email" xml:"email" validate:"email"`


}






// UserListItem


// Username 用户名

// Email 邮箱

// Age 年龄

// CreatedAt 创建时间

// UpdatedAt 更新时间

// Id 用户ID

type UserListItem struct {

	Username string `json:"username" validate:"required,min=3,max=20"`

	Email string `json:"email" validate:"required,email"`

	Age int `json:"age" validate:"min=0,max=150"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	Id int64 `json:"id" validate:"required"`

}

// UserData


// Id 用户ID

// Username 用户名

// Email 邮箱

// Age 年龄

// CreatedAt 创建时间

// UpdatedAt 更新时间

type UserData struct {

	Id int64 `json:"id" validate:"required"`

	Username string `json:"username" validate:"required,min=3,max=20"`

	Email string `json:"email" validate:"required,email"`

	Age int `json:"age" validate:"min=0,max=150"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

}

