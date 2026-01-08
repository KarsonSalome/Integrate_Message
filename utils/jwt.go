package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("dev-secret")

func GenerateToken(uid uint) (string, error) {
    claims := jwt.MapClaims{
        "uid": uid,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
