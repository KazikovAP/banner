package jwt

import (
	logerr "banner/internal/lib/logger/logerr"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Require struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type JWTSecret struct {
	secret []byte
	log    *slog.Logger
}

func NewJWTSecret(secret string, log *slog.Logger) *JWTSecret {
	return &JWTSecret{secret: []byte(secret), log: log}
}

func (secret *JWTSecret) GenerateToken(username string, role string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret.secret)
	if err != nil {
		secret.log.Error("Failed to sign token")
		return "", fmt.Errorf("failed to sign token", logerr.Err(err))
	}

	return signedToken, nil
}

func (secret *JWTSecret) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			secret.log.Error("Unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}

		return secret.secret, nil
	})
	if err != nil {
		secret.log.Error("Failed to parse token", logerr.Err(err))
		return nil, fmt.Errorf("failed to parse token", logerr.Err(err))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		secret.log.Error("Invalid token")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (secret *JWTSecret) ExtractRoleFromToken(tokenString string) (string, error) {
	claims, err := secret.VerifyToken(tokenString)
	if err != nil {
		secret.log.Error("Error with token", logerr.Err(err))
		return "", err
	}

	role, ok := claims["role"].(string)
	if !ok {
		secret.log.Error("Role not found in claims")
		return "", errors.New("role not found in token")
	}

	return role, nil
}

func (secret *JWTSecret) ExtractRoleAndUsernameFromToken(tokenString string) (string, string, error) {
	claims := &Require{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			secret.log.Error("Unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}

		return secret.secret, nil
	})

	if err != nil {
		secret.log.Error("Failed to parse token", logerr.Err(err))
		return "", "", fmt.Errorf("failed to parse token", logerr.Err(err))
	}

	if !token.Valid {
		secret.log.Error("Invalid token")
		return "", "", errors.New("invalid token")
	}

	return claims.Username, claims.Role, nil
}
