/*
 * Generated from OpenAPI3 document
 * Module: User
 */

package param





// UserListUsersRequest 获取用户列表
type UserListUsersRequest struct {

	Page int `json:"page" validate:"min=1"` // 页码

	Size int `json:"size" validate:"min=1,max=100"` // 每页数量


}



// UserCreateRequest 创建用户
type UserCreateRequest struct {



	Age int `json:"age" validate:"min=0,max=150"` // 

	Username string `json:"username" validate:"required" validate:"min=3,max=20"` // 

	Email string `json:"email" validate:"required" validate:"email"` // 


}




// UserListItem 
type UserListItem struct {

	Username string `json:"username"` // 

	Email string `json:"email"` // 

	Id int `json:"id"` // 

}

// UserData 
type UserData struct {

	Id int `json:"id"` // 

	Username string `json:"username"` // 

	Email string `json:"email"` // 

}

