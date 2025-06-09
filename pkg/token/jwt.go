package token

import (
	"fmt"
	"time"

	"github.com/EduardoMark/my-finance-api/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret string
}

func NewTokenManager(cfg config.Env) *TokenManager {
	return &TokenManager{secret: cfg.JWTSecret}
}

func (s *TokenManager) GenerateToken(name, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"name":  name,
			"email": email,
			"exp":   time.Now().Add(time.Hour * 1).Unix(),
		})

	tokenStr, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("error on parse token to string: %w", err)
	}

	return tokenStr, nil
}

func (s *TokenManager) VerifyToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
