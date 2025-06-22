package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

func GenerateToken(email string, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {

	if !strings.HasPrefix(token, "Bearer ") {
		return 0, errors.New("authorization header invalid")
	}
	token = strings.TrimPrefix(token, "Bearer ")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userId, ok := claims["id"].(float64)
	fmt.Println("User ID from claims:", userId, claims)
	if !ok {
		return 0, errors.New("something went wrong")
	}

	// userId := claims["userId"]

	return int64(userId), nil

}
