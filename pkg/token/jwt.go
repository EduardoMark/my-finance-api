package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/EduardoMark/my-finance-api/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret string
}

func NewTokenManager(cfg config.Env) *TokenManager {
	return &TokenManager{secret: cfg.JWTSecret}
}

func (s *TokenManager) GenerateToken(name, email string) (string, error) {
	claims := Claims{
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("error on parse token to string: %w", err)
	}

	return tokenStr, nil
}

func (s *TokenManager) VerifyToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
