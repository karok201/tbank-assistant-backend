package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var jwtSecret = []byte(viper.GetString("JWT_SECRET"))

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": strconv.Itoa(userID),
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Токен истекает через 24 часа
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func GetUserIdFromToken(c *gin.Context) (int, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return -1, errors.New("invalid token")
	}

	if len(authHeader) < 7 || authHeader[:7] != "ApiKey " {
		return -1, errors.New("invalid token")
	}

	tokenString := authHeader[7:]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return -1, errors.New("invalid token")
	}

	// Проверяем валидность токена
	if !token.Valid {
		return -1, errors.New("invalid token")
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		return -1, fmt.Errorf("invalid user ID: %v", err)
	}

	return userID, nil
}
