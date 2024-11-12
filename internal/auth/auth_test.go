package auth

import (
	"errors"
	"net/http"
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
	t.Error("Token did not expire")
}

func TestMalformedToken(t *testing.T) {
	_, err := ValidateJWT("random.string_not_valid.jwt", password)
	if err == nil {
		t.Error("Expected an error for malformed token, got nil")
	}
	if !errors.Is(err, jwt.ErrTokenMalformed) {
		t.Errorf("Expected an error for malformed token, got: %v\n", err)
	}

	_, err = ValidateJWT("", password)
	if err == nil {
		t.Error("Expected an error for empty token, got: nil")
	}
	if !errors.Is(err, jwt.ErrTokenMalformed) {
		t.Errorf("Expected an error for malformed token, got: %v\n", err)
	}
}

func TestBearerToken(t *testing.T) {
	expected := "1235234lkjsadfasd"
	header := make(http.Header)
	header.Set("Authorization", "Bearer "+expected)

	actual, err := GetBearerToken(header)
	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Errorf("Expected bearer token to equal %v, got %v\n", expected, actual)
	}
}

func TestMalformedBearerToken(t *testing.T) {
	malformedError := "malformed header"
	missingAutherror := "missing authorization header"
	header := make(http.Header)
	header.Set("Authorization", "Bearer ")
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatal("expected malformed token error, got nil")
	} else {
		if err.Error() != malformedError {
			t.Errorf("expected error message to equal '%v', got %v\n", malformedError, err)
		}
	}

	header2 := make(http.Header)
	_, err = GetBearerToken(header2)
	if err == nil {
		t.Fatal("expected malformed token error, got nil")
	} else {
		if err.Error() != missingAutherror {
			t.Errorf("expected error message to equal '%v', got %v\n", missingAutherror, err)
		}
	}

}
