/*
 *
 * user.go
 * models
 *
 * Created by lintao on 2019-01-29 16:31
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package models

type User struct {
	Id       int64  `json:"id" form:"id" query:"id" json:"id" xorm:"not null pk autoincr INT(10)"`
	Name     string `json:"name" form:"name" query:"name" xorm:"not null default '' VARCHAR(32)"`
	Phone    string `json:"phone" form:"phone" query:"phone" xorm:"not null default '' unique VARCHAR(64)"`
	Status   int64  `json:"status" form:"status" query:"status"`
	Account  string `json:"account" form:"account" query:"account" xorm:"not null default '' index VARCHAR(16)"`
	Password string `json:"password" form:"password" query:"password"`
	Created  Time   `json:"created" form:"created" query:"created"`
}

type UserParam struct {
	ApiQuery
	User
}
