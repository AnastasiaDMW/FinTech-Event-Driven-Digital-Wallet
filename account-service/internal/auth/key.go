package auth

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func LoadPublicKey() (*rsa.PublicKey, error) {
	pub := os.Getenv("PUBLIC_KEY")

	if pub == "" {
		return nil, ErrMissingKeys
	}

	keyPublic, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pub))
	if err != nil {
		return nil, err
	}

	return keyPublic, nil
}
