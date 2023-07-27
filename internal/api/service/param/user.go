/*
 * Created by lintao on 2023/7/26 下午2:39
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package param

import (
	"github.com/NSObjects/go-template/internal/api/data/model"
)

type UserParam struct {
	APIQuery
	model.User
}

type UserResponse struct {
	Name     string `json:"name" form:"name" query:"name"`
	Phone    string `json:"phone" form:"phone" query:"phone"`
	Status   int64  `json:"status" form:"status" query:"status"`
	Password string `json:"password" form:"password" query:"password"`
}
