package auth

import (
	"l3/WarehouseControl/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secretKey string
	ttl       time.Duration
}

func NewJWT(secretKey string, ttl time.Duration) *JWT {
	return &JWT{
		secretKey: secretKey,
		ttl:       ttl,
	}
}

func (j *JWT) GenerateToken(user *models.User) (string, error) {
	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (j *JWT) ParseToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	//TODO подумать о валидности токена
	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(j.secretKey), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
