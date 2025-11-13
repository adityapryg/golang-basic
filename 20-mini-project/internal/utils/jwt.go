package utils

import (
	"errors"
	"time"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims struktur JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT token untuk user
func GenerateToken(userID uint, username string) (string, error) {
	// Token berlaku 24 jam
	expirationTime := time.Now().Add(24 * time.Hour)

	// Buat claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key dari env
	cfg := config.LoadConfig()
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memvalidasi JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parse token
	cfg := config.LoadConfig()
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
