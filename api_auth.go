package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aramirez3/chirpy/internal/auth"
	"github.com/aramirez3/chirpy/internal/database"
)

type RefreshTokenResponse struct {
	Token string
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

func (cfg *apiConfig) createJWT(ctx context.Context, user *User) error {
	fmt.Println("Create new JWT for user")
	// create the jwt
	// save it in the db with updated timestamps
	// dont need to return value
	newToken
}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handle refresh")
	reqBody, err := io.ReadAll(req.Body)
	if err != nil || len(reqBody) > 0 {
		fmt.Printf("err value: %v, isNil %v\n", err, err == nil)
		fmt.Printf("reqBody: %v, hasLen %v\n", reqBody, len(reqBody) > 0)
		returnBadRequest(w)
		return
	}
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		returnNotAuthorized(w)
		return
	}
	if refreshToken == "" {
		fmt.Println("bear token is empty")
		returnNotAuthorized(w)
		return
	}
	dbToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		fmt.Printf("error making db call: %v\n", err)
		returnNotAuthorized(w)
		return
	}

	if dbToken.ExpiresAt.Before(time.Now().UTC()) {
		fmt.Println("refresh token is expired")
		returnNotAuthorized(w)
		return
	} else {
		fmt.Println("refresh token is still valid")
	}
	// need the user
	// refresh the jwt for the user
	newToken := RefreshTokenResponse{
		// should be the refreshed token from user
		Token: dbToken.Token,
	}
	respBody, err := encodeJson(newToken)

	if err != nil {
		fmt.Printf("error encoding token response %v\n", err)
		returnErrorResponse(w, standardError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}
