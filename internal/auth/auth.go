package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userId.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}

	sub := claims.Subject
	if sub == "" {
		return uuid.UUID{}, fmt.Errorf("subject is empty")
	}

	userId, err := uuid.Parse(sub)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("missing authorization header")
	}
	words := strings.Split(tokenString, " ")
	if words[0] != "Bearer" || words[1] == "" {
		return "", fmt.Errorf("malformed header")
	}
	return words[1], nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	bytes.Equal(b, make([]byte, 32))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	if auth == "" {
		return "", fmt.Errorf("no Authorization key provided")
	}

	authValues := strings.Split(auth, " ")
	if authValues[0] != "ApiKey" || authValues[1] == "" {
		return "", fmt.Errorf("invalid ApiKey")
	}
	return authValues[1], nil
}
