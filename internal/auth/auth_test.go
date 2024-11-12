package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHash(t *testing.T) {
	password := "password123"
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
	password := "password123"
	wronghash, _ := HashPassword("helloworld")

	err := CheckPasswordHash(password, wronghash)
	if err == nil {
		t.Fatal(err)
	}
}

func TestJWT(t *testing.T) {
	expected := uuid.New()
	tokenString, err := MakeJWT(expected, "secret", 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ValidateJWT(tokenString, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Errorf("Invalid UUID: expected %v, actual %v\n", expected, actual)
	}
}

func TestJWTMismatch(t *testing.T) {
	expected := uuid.New()
	tokenString, err := MakeJWT(expected, "secret", 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ValidateJWT(tokenString, "NoSecret")
	if err == nil {
		t.Fatal(err)
	}
	if expected == actual {
		t.Errorf("UUIDs should not match: expected %v, actual %v\n", expected, actual)
	}

}
