package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	password      = "password123"
	wrongPassword = "grassword123"
)

func TestPasswordHash(t *testing.T) {
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPasswordMismatch(t *testing.T) {
	wronghash, _ := HashPassword(wrongPassword)

	err := CheckPasswordHash(password, wronghash)
	if err == nil {
		t.Fatal(err)
	}
}

func TestJWT(t *testing.T) {
	expected := uuid.New()
	tokenString, err := MakeJWT(expected, password, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ValidateJWT(tokenString, password)
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Errorf("Invalid UUID: expected %v, actual %v\n", expected, actual)
	}
}

func TestJWTMismatch(t *testing.T) {
	expected := uuid.New()
	tokenString, err := MakeJWT(expected, password, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ValidateJWT(tokenString, wrongPassword)
	if err == nil {
		t.Fatal(err)
	}
	if expected == actual {
		t.Errorf("UUIDs should not match: expected %v, actual %v\n", expected, actual)
	}

}

func TestExpiredToken(t *testing.T) {
	id := uuid.New()
	duration := time.Duration(2 * time.Second)
	tokenString, err := MakeJWT(id, password, duration)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(duration * 2)
	_, err = ValidateJWT(tokenString, password)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			t.Log("Token expired as expected")
		} else {
			t.Errorf("Expected token expiration error, got: %v\n", err)
		}
		return
	}
	t.Errorf("Token did not expire")
}
