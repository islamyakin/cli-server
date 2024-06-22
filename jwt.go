package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	secretKey = []byte("secret -key")
)

func CreateToken(username, email, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"email":    email,
			"surname":  name,
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(time.Minute * 30).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
