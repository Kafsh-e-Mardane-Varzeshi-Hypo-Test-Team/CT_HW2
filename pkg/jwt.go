package pkg

import (
	"errors"
	"time"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	SessionMaxAge = 24 * time.Hour
)

func GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(SessionMaxAge).Unix(),
	})

	return token.SignedString(config.JWT.SecretKey)
}

func ValidateToken(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// safety measure
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return config.JWT.SecretKey, nil
	})

	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, false
	}

	// Check if the token is expired
	if claims["exp"] == nil {
		return nil, false
	}

	exp, isFloat := claims["exp"].(float64)
	if !isFloat {
		return nil, false
	}
	if time.Now().Unix() > int64(exp) {
		return nil, false
	}

	// Check if the userId is present in the claims
	if claims["userId"] == nil {
		return nil, false
	}

	return claims, true
}
