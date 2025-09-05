/*
 * Generated test cases from OpenAPI3 document
 * Module: User
 */

package biz

import (
	"context"
	"testing"

	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)



func TestUserHandler_ListUsers(t *testing.T) {
	// 创建handler
	handler := &UserHandler{
		dataManager: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.UserListUsersRequest{
		Page:  1,
		Size: 10,
	}

	// 测试ListUsers方法
	result, total, err := handler.ListUsers(ctx, req)

	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	assert.NoError(t, err)
}

func TestUserHandler_ListUsers_Validation(t *testing.T) {
	// 测试参数验证
	validator := validator.New()
	
	// 测试有效数据
	validReq := param.UserListUsersRequest{
		Page:  1,
		Size: 10,
	}
	err := validator.Struct(validReq)
	assert.NoError(t, err, "有效数据应该通过验证")
	
	// 测试无效数据
	invalidReq := param.UserListUsersRequest{
		Page:  -1, // 无效页码
		Size:  -1, // 无效页面大小
	}
	err = validator.Struct(invalidReq)
	assert.Error(t, err, "无效数据应该验证失败")
}



func TestUserHandler_Create(t *testing.T) {
	// 创建handler
	handler := &UserHandler{
		dataManager: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.UserCreateRequest{}

	// 测试Create方法
	err := handler.Create(ctx, req)

	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.NoError(t, err)
}

func TestUserHandler_Create_Validation(t *testing.T) {
	// 测试参数验证
	validator := validator.New()
	
	// 测试有效数据
	validReq := param.UserCreateRequest{


		Age: 1, // 

		Username: "测试Username", // 

		Email: "测试Email", // 


	}
	err := validator.Struct(validReq)
	assert.NoError(t, err, "有效数据应该通过验证")
	
	// 测试无效数据
	invalidReq := param.UserCreateRequest{
		// 空结构体，所有必填字段都缺失
	}
	err = validator.Struct(invalidReq)
	assert.Error(t, err, "无效数据应该验证失败")
}


