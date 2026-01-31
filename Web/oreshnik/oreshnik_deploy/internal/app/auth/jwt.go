package auth

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func LoadKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privBytes, err := os.ReadFile("public/private.pem")
	if err != nil {
		return nil, nil, fmt.Errorf("read private key: %w", err)
	}
	priv, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse private key: %w", err)
	}

	pubBytes, err := os.ReadFile("public/public.pem")
	if err != nil {
		return nil, nil, fmt.Errorf("read public key: %w", err)
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse public key: %w", err)
	}
	return priv, pub, nil
}

func GenerateJWT(userID uint, isAdmin bool, privateKey *rsa.PrivateKey) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func Parse(tokenString string, publicKey *rsa.PublicKey) (*jwt.Token, jwt.MapClaims, error) {
	tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}
	return tok, claims, nil
}
