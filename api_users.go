package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aramirez3/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type CreateUser struct {
	Email string
}

func (cfg *apiConfig) handleNewUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	createUser := CreateUser{}

	w.Header().Add(contentType, plainTextContentType)

	err := decoder.Decode(&createUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		returnErrorResponse(w)
		return
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     createUser.Email,
	}
	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), params)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		returnErrorResponse(w)
		return
	}
	if dbUser.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		returnErrorResponse(w)
		return
	}

	newUser := User{
		Id:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respBody, _ := encodeJson(newUser)
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}
