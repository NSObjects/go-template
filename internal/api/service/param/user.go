/*
 * Generated from OpenAPI3 document
 * Module: User
 */

package param

import (
	"time"
)



// UserListUsersRequest 获取用户列表
type UserListUsersRequest struct {

	Page int `json:"page"` // 页码

	Size int `json:"size"` // 每页数量


}



// UserCreateRequest 创建用户
type UserCreateRequest struct {



	Username string `json:"username" validate:"required" validate:"min=3,max=20"` // 用户名

	Email string `json:"email" validate:"required" validate:"email"` // 邮箱

	Age int `json:"age" validate:"min=0,max=150"` // 年龄


}





// UserUpdateRequest 更新用户
type UserUpdateRequest struct {



	Age int `json:"age" validate:"min=0,max=150"` // 年龄

	Username string `json:"username" validate:"min=3,max=20"` // 用户名

	Email string `json:"email" validate:"email"` // 邮箱


}






// UserListItem 
type UserListItem struct {

	CreatedAt time.Time `json:"created_at"` // 创建时间

	UpdatedAt time.Time `json:"updated_at"` // 更新时间

	Id int64 `json:"id"` // 用户ID

	Username string `json:"username"` // 用户名

	Email string `json:"email"` // 邮箱

	Age int `json:"age"` // 年龄

}

// UserData 
type UserData struct {

	UpdatedAt time.Time `json:"updated_at"` // 更新时间

	Id int64 `json:"id"` // 用户ID

	Username string `json:"username"` // 用户名

	Email string `json:"email"` // 邮箱

	Age int `json:"age"` // 年龄

	CreatedAt time.Time `json:"created_at"` // 创建时间

}

