package auth

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

type Claims struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GenerateTokenPair(userID, email, role string) (*TokenPair, error) {
	accessToken, err := generateToken(userID, email, role, "access", 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(userID, email, role, "refresh", 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ValidateToken(tokenString, expectedType string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			_, ok := token.Header["alg"].(string)
			if !ok {
				return nil, ErrInvalidAlg
			}
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != expectedType {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

func InitKeys() error {
	keyPrivData, err := os.ReadFile("private.pem")
	if err != nil {
		return err
	}
	keyPubData, err := os.ReadFile("public.pem")
	if err != nil {
		return err
	}

	keyPrivate, err := jwt.ParseRSAPrivateKeyFromPEM(keyPrivData)
	if err != nil {
		return err
	}
	keyPublic, err := jwt.ParseRSAPublicKeyFromPEM(keyPubData)
	if err != nil {
		return err
	}

	privateKey = keyPrivate
	publicKey = keyPublic
	return nil
}

func generateToken(userID, email, role, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}
