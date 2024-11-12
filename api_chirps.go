package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aramirez3/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirps struct {
	Chirps []Chirp
}

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type ChirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidResponse struct {
	Valid bool `json:"valid"`
}

func (cfg *apiConfig) handleNewChirp(w http.ResponseWriter, req *http.Request) {
	reqChirp := ChirpRequest{}
	isValid, errorString := validateChirpRequest(req.Body, &reqChirp)
	w.Header().Add(contentType, plainTextContentType)
	if errorString != "" || !isValid {
		returnErrorResponse(w, errorString)
		return
	}

	chirp := Chirp{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      reqChirp.Body,
		UserId:    reqChirp.UserId,
	}
	encodedChirp, err := encodeJson(chirp)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	params := database.CreateChirpParams{
		ID:        chirp.Id,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserId,
	}
	_, err = cfg.dbQueries.CreateChirp(req.Context(), params)

	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(encodedChirp)
}

func validateChirpRequest(body io.ReadCloser, chirp *ChirpRequest) (bool, string) {
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

func removeProfanity(chirp *ChirpRequest) {
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

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	responseChirps := []Chirp{}
	if len(dbChirps) > 0 {
		responseChirps = dbChirpsToResponse(dbChirps)
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseChirps)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	w.Header().Add(contentType, plainTextContentType)
}

func dbChirpsToResponse(dbChirps []database.Chirp) []Chirp {
	response := []Chirp{}
	for _, c := range dbChirps {
		response = append(response, dbChirpToResponse(c))
	}
	return response
}

func dbChirpToResponse(c database.Chirp) Chirp {
	return Chirp{
		Id:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserId:    c.UserID,
	}
}
