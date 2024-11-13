package main

import (
	"context"
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
}

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login = CreateUser

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

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	login := Login{}
	err := decoder.Decode(&login)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(req.Context(), login.Email)
	if err != nil || dbUser.HashedPassword == "" {
		returnErrorResponse(w, standardError)
		return
	}

	err = auth.CheckPasswordHash(login.Password, dbUser.HashedPassword)
	if err != nil {
		returnNotAuthorized(w)
		return
	}

	jwt, err := auth.MakeJWT(dbUser.ID, cfg.Secret, time.Hour)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	w.WriteHeader(http.StatusOK)
	user := ToResponseUser(dbUser)
	err = cfg.createRefreshToken(req.Context(), &user)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	user.Token = jwt

	responseUser, err := encodeJson(user)
	if err != nil {
		returnNotAuthorized(w)
		return
	}
	w.Header().Set(contentType, plainTextContentType)
	w.Write(responseUser)
}

func ToResponseUser(u database.User) User {
	return User{
		Id:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}

func (cfg *apiConfig) createRefreshToken(ctx context.Context, user *User) error {
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return err
	}
	params := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.Id,
		ExpiresAt: time.Now().UTC().Add(24 * 60 * time.Hour),
	}

	_, err = cfg.dbQueries.CreateRefreshToken(ctx, params)
	if err != nil {
		return err
	}

	user.RefreshToken = refreshToken
	return nil
}
