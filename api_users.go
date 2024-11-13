package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aramirez3/chirpy/internal/auth"
	"github.com/aramirez3/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type CreateUser = UserPreferences

type Login = UserPreferences

type UserPreferences struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handleNewUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	createUser := CreateUser{}

	w.Header().Add(contentType, plainTextContentType)

	err := decoder.Decode(&createUser)

	if err != nil || createUser.Email == "" || createUser.Password == "" {
		returnErrorResponse(w, standardError)
		return
	}

	hash, err := auth.HashPassword(createUser.Password)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	params := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Email:          createUser.Email,
		HashedPassword: hash,
	}
	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), params)

	if err != nil || dbUser.Email == "" {
		returnErrorResponse(w, standardError)
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

func (cfg *apiConfig) handleUserUpdate(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnNotAuthorized(w)
		return
	}

	jwtId, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		if err.Error() == "invalid token" || err.Error() == "subject is empty" {
			returnBadRequest(w)
		} else {
			returnNotAuthorized(w)
		}
		return
	}

	payload := UserPreferences{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	hash, err := auth.HashPassword(payload.Password)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	params := database.UpdateUserParams{
		ID:             jwtId,
		Email:          payload.Email,
		UpdatedAt:      time.Now().UTC(),
		HashedPassword: hash,
	}

	dbUser, err := cfg.dbQueries.UpdateUser(req.Context(), params)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	updatedUser := User{
		Id:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respBody, _ := encodeJson(updatedUser)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
