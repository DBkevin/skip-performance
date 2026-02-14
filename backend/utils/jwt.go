package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"skin-performance/models"
)

const jwtSecret = "skin-performance-secret-key-2024"
const tokenExpireDuration = time.Hour * 24 * 7 // 7天

type Claims struct {
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	EmployeeID *uint  `json:"employee_id,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID:     user.ID,
		Username:   user.Username,
		Role:       user.Role,
		EmployeeID: user.EmployeeID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "skin-performance",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
