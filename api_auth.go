package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/aramirez3/chirpy/internal/auth"
	"github.com/aramirez3/chirpy/internal/database"
)

type RefreshTokenResponse struct {
	Token string `json:"token"`
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

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil || len(reqBody) > 0 {
		returnBadRequest(w)
		return
	}
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnNotAuthorized(w)
		return
	}
	if refreshToken == "" {
		returnNotAuthorized(w)
		return
	}
	dbToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		returnNotAuthorized(w)
		return
	}

	if dbToken.ExpiresAt.Before(time.Now().UTC()) {
		returnNotAuthorized(w)
		return
	} else {
	}
	jwt, err := auth.MakeJWT(dbToken.UserID, cfg.Secret, time.Hour)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	newToken := RefreshTokenResponse{
		Token: jwt,
	}
	respBody, err := encodeJson(newToken)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	w.WriteHeader(http.StatusOK)
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

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil || len(reqBody) > 0 {
		returnBadRequest(w)
		return
	}
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnNotAuthorized(w)
		return
	}
	if refreshToken == "" {
		returnNotAuthorized(w)
		return
	}
	now := time.Now().UTC()
	params := database.UpdateRefreshTokenParams{
		Token:     refreshToken,
		UpdatedAt: now,
		RevokedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}
	_, err = cfg.dbQueries.UpdateRefreshToken(req.Context(), params)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
