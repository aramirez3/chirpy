package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aramirez3/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string `json:"body"`
	UserId    uuid.UUID
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidResponse struct {
	Valid bool `json:"valid"`
}

func (cfg *apiConfig) createTestUser(req *http.Request) database.User {
	createUser := CreateUser{}
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     createUser.Email,
	}
	dbUser, _ := cfg.dbQueries.CreateUser(req.Context(), params)
	return dbUser
}

func handleNewChirp(w http.ResponseWriter, req *http.Request) {
	user := cfg.createTestUser(req)
	chirp := Chirp{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	isValid, errorString := validateChirp(req.Body, &chirp)
	w.Header().Add(contentType, plainTextContentType)
	if errorString != "" || !isValid {
		returnErrorResponse(w, errorString)
		return
	}

	encodedChirp, err := encodeJson(chirp)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	log.Printf("This chirp gets encoded:\n%v\n", chirp)
	w.WriteHeader(http.StatusCreated)
	w.Write(encodedChirp)
}

func validateChirp(body io.ReadCloser, chirp *Chirp) (bool, string) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&chirp)
	if err != nil {
		return false, standardError
	}
	if chirp.Body == "" {
		return false, standardError
	}

	if len(chirp.Body) > 140 {
		return false, "Chirp is too long"
	}

	removeProfanity(chirp)
	return true, ""
}

func removeProfanity(chirp *Chirp) {
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	words := strings.Split(chirp.Body, " ")
	if len(words) > 0 {
		for i, word := range words {
			w, ok := badWords[strings.ToLower(word)]
			if w && ok {
				if word != "Sharbert!" {
					words[i] = "****"
				}
			}
		}
		chirp.Body = strings.Join(words, " ")
	}
}
