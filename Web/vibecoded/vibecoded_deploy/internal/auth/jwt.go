package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"vibecoded/internal/models"
)

const (
	BearerTokenPrefix = "Bearer "
	tokenExpiration   = 10 * time.Minute
)

type JWTProvider struct {
	jwtSecret []byte
}

func NewJWTProvider(jwtSecret []byte) *JWTProvider {
	return &JWTProvider{
		jwtSecret: jwtSecret,
	}
}

func (j *JWTProvider) NewJWTToken(user *models.UserDBItem) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      strconv.Itoa(user.ID),
		"username": user.Username,
		"isAdmin":  user.IsAdmin,
		"exp":      time.Now().Add(tokenExpiration).Unix(),
		"iat":      time.Now().Unix(),
	})

	token, err := claims.SignedString(j.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("claims.SignedString: %w", err)
	}

	return token, nil
}

func (j *JWTProvider) verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (any, error) {
			return j.jwtSecret, nil
		})
	if err != nil {
		return nil, fmt.Errorf("jwt.Parse: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}

func (j *JWTProvider) ParseJWTToken(r *http.Request) (*jwt.Token, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Value == "" {
		return nil, ErrMissingAuthCookie
	}
	if !strings.HasPrefix(tokenCookie.Value, BearerTokenPrefix) {
		return nil, ErrAuthCookieNotBearer
	}
	tokenString := strings.TrimPrefix(tokenCookie.Value, BearerTokenPrefix)

	token, err := j.verifyToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("token failed to verify: %w", err)
	}

	return token, nil
}

func (j *JWTProvider) SetJWTTokenCookie(w http.ResponseWriter, tokenString string,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(tokenExpiration),
		Path:     "/",
		HttpOnly: true,
	})
}
