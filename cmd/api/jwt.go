package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (app *application) generateToken(userID int64, email string) (string, error) {
	expiryDuration, err := time.ParseDuration(app.config.jwtConfig.expiry)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.Itoa(int(userID)),
			Issuer:    app.config.jwtConfig.issuer,
			Audience:  []string{email},
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
