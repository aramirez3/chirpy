package auth

import "testing"

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
