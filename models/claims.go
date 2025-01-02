package models

import "github.com/dgrijalva/jwt-go"

// Claims структура для хранения данных из JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
