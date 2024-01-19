/*
 *
 * jwt.go
 * data
 *
 * Created by lintao on 2023/11/21 15:21
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package data

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	Name  string `json:"name"`
	ID    int64  `json:"id" `
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}
