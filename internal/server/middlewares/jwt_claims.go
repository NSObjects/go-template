package middlewares

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	Name  string `json:"name"`
	ID    int64  `json:"id" `
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}
