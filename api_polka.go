package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aramirez3/chirpy/internal/auth"
	"github.com/aramirez3/chirpy/internal/database"
	"github.com/google/uuid"
)

type PolkaRequest struct {
	Event string           `json:"event"`
	Data  PolkaRequestData `json:"data"`
}

type PolkaRequestData struct {
	UserId uuid.UUID `json:"user_id"`
}

const (
	validPolkaWebhookEvent = "user.upgraded"
)

func (cfg *apiConfig) handlePolkaWebhooks(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil || apiKey != cfg.PolkaKey {
		returnUnauthorized(w)
		return
	}

	decoder := json.NewDecoder(req.Body)
	requestData := PolkaRequest{}
	err = decoder.Decode(&requestData)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	if requestData.Event != validPolkaWebhookEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	params := database.UpgradeUserToRedParams{
		ID: requestData.Data.UserId,
		IsChirpyRed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		UpdatedAt: time.Now().UTC(),
	}

	dbUser, err := cfg.dbQueries.GetUserById(req.Context(), requestData.Data.UserId)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	if dbUser.IsChirpyRed.Valid && !dbUser.IsChirpyRed.Bool {
		_, err := cfg.dbQueries.UpgradeUserToRed(req.Context(), params)
		if err != nil {
			returnErrorResponse(w, standardError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
