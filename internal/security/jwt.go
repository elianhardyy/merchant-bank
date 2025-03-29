package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(username, email string, role ...string) (string, error)
	IsTokenExpired(tokenString string) bool
	VerifyToken(tokenString string) (*Claims, error)
	IsTokenValid(tokenString string) bool
}

type tokenService struct {
	jwtSecret []byte
}

func NewTokenService(secret []byte) TokenService {
	return &tokenService{jwtSecret: secret}
}

type Claims struct {
	Email string   `json:"email"`
	Role  []string `json:"role"`
	jwt.RegisteredClaims
}

func (t *tokenService) GenerateToken(username string, email string, role ...string) (string, error) {
	claims := &Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			Issuer:    "go-json",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.jwtSecret)
}

func (t *tokenService) IsTokenExpired(tokenString string) bool {
	claims, err := t.VerifyToken(string(t.jwtSecret))
	if err != nil {
		return true
	}
	return claims.ExpiresAt.Unix() < time.Now().Unix()
}

func (t *tokenService) IsTokenValid(tokenString string) bool {
	_, err := t.VerifyToken(string(t.jwtSecret))
	return err == nil
}

func (t *tokenService) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return token.Claims.(*Claims), nil
}
