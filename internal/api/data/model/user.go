/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" form:"name" query:"name" binding:"required"`
	Phone    string `json:"phone" form:"phone" query:"phone"`
	Status   int64  `json:"status" form:"status" query:"status"`
	Password string `json:"password" form:"password" query:"password"`
}
