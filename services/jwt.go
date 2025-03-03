package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) GetSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "adr_keren"
	}
	return secretKey
}

func (s *JWTService) ExtractToken(tokenString string) string {
	if strings.HasPrefix(tokenString, "Bearer ") {
		return strings.Split(tokenString, " ")[1]
	}
	return tokenString
}

func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(s.GetSecretKey()), nil
	})

	return token, err
}

func (s *JWTService) ExtractUserID(token *jwt.Token) (uint, error) {
	claims := token.Claims.(jwt.MapClaims)
	var userID uint

	switch v := claims["user_id"].(type) {
	case string:
		var id int
		_, err := fmt.Sscanf(v, "%d", &id)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id format")
		}
		userID = uint(id)
	case float64:
		userID = uint(v)
	default:
		return 0, fmt.Errorf("unexpected type for user_id: %T", v)
	}

	return userID, nil
}