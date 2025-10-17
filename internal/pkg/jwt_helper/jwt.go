package jwt_helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTHelper struct {
	secret []byte
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJwtHelper(secret string) (*JWTHelper, error) {
	if secret == "" {
		return nil, errors.New("empty JWT secret")
	}
	return &JWTHelper{secret: []byte(secret)}, nil
}

func (h *JWTHelper) GenerateJWT(userID, email string, expiresIn time.Duration) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "meetly-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secret)
}

func (h *JWTHelper) GenerateDefaultToken(userID, email string) (string, error) {
	return h.GenerateJWT(userID, email, 24*time.Hour)
}

func (h *JWTHelper) ParseJWT(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token string")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return h.secret, nil
		},
	)

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (h *JWTHelper) ValidateToken(tokenString string) (bool, error) {
	_, err := h.ParseJWT(tokenString)
	if err != nil {
		return false, err
	}
	return true, nil
}
