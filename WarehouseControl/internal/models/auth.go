package models

import "github.com/golang-jwt/jwt/v5"

// модель авторизованного юзера
type CurrentUser struct {
	UserID   int64
	Username string
	Role     string
}

// структура внутри JWT
type Claims struct {
	UserID   int64
	Username string
	Role     string
	jwt.RegisteredClaims
}
