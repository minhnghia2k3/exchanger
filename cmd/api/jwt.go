package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"time"
)

type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func (app *application) generateToken(user *store.User) (string, error) {
	expiryDuration, err := time.ParseDuration(app.config.jwtConfig.expiry)
	if err != nil {
		return "", err
	}

	identifier, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        identifier.String(),
			Issuer:    app.config.jwtConfig.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
		},
	})

	tokenString, err := token.SignedString([]byte(app.config.jwtConfig.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (app *application) generateRefreshToken(userID int64) (string, error) {
	expiryDuration, err := time.ParseDuration(app.config.jwtConfig.refreshExpiry)
	if err != nil {
		return "", err
	}

	identifier, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        identifier.String(),
			Issuer:    app.config.jwtConfig.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(app.config.jwtConfig.secret))
}

func (app *application) verifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.config.jwtConfig.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		return nil, fmt.Errorf("invalid token")
	} else {
		return claims, nil
	}
}
