package model

import "github.com/dgrijalva/jwt-go"

type JwtCustomClaims struct {
	Mail  string `json:"mail"`
	jwt.StandardClaims
}
