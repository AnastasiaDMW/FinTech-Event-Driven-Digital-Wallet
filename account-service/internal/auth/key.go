package auth

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func LoadPublicKey(path string) *rsa.PublicKey {
	keyData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		panic(err)
	}

	return publicKey
}