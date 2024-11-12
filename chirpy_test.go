// calculator_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestReadiness(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleReadiness)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handleReadines.StatusCode: actual %v, expected %v\n", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handleReadines.Body: actual %v, expected %v\n", rr.Body.String(), expected)
	}
}

func TestNewChirp(t *testing.T) {
	chirp := Chirp{
		Body:   "Hello world",
		UserId: uuid.New(),
	}

	jsonData, err := json.Marshal(chirp)
	if err != nil {
		t.Fatal(err)
	}

	reqBody := bytes.NewReader(jsonData)
	req, err := http.NewRequest("POST", "/api/chirps", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleNewChirp)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handleReadines.StatusCode: actual %v, expected %v\n", status, http.StatusCreated)
	}

	responseChirp := Chirp{}

	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&responseChirp)
	if err != nil {
		t.Fatal(err)
	}

	if responseChirp.Body != chirp.Body {
		t.Errorf("Invalid body: expected %v, actual %v\n", responseChirp.Body, chirp.Body)
	}
	if responseChirp.CreatedAt != chirp.CreatedAt {
		t.Errorf("Invalid created_at: expected %v, actual %v\n", responseChirp.CreatedAt, chirp.CreatedAt)
	}
	if responseChirp.UpdatedAt != chirp.UpdatedAt {
		t.Errorf("Invalid updated_at: expected %v, actual %v\n", responseChirp.UpdatedAt, chirp.UpdatedAt)
	}
	if responseChirp.UserId != chirp.UserId {
		t.Errorf("Invalid user_id: expected %v, actual %v\n", responseChirp.UserId, chirp.UserId)
	}
}
