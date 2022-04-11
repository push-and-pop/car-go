package token

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenToken(phone string) {
	claim := jwt.StandardClaims{
		Issuer:   phone,
		IssuedAt: time.Now().Unix(),
	}
}
