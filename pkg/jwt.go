package pkg

import (
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: function that loads all env variables, log.Fatals any error if unsuccessful
// prolly add a map and function for retrieving env variables (no error output)
var js, _ = os.LookupEnv("JWT_SECRET")
var jwtSecret, _ = base64.StdEncoding.DecodeString(js)

const (
	sessionMaxAge = 24 * time.Hour
)

func GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(sessionMaxAge).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// safety measure
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, false
	}

	return claims, true
}
