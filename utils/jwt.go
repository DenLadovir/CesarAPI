package utils

import (
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"time"
)

var jwtKey = []byte("losharanekantara")

type Claims struct {
	Username string `json: "username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		log.Printf("JWT Validation error: %v", err)
		return nil, err
	}
	return claims, nil
}
