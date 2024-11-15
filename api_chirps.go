package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aramirez3/chirpy/internal/auth"
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
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnUnauthorized(w)
		return
	}

	jwtId, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		returnUnauthorized(w)
		return
	}

	reqChirp := ChirpRequest{}
	isValid, errorString := validateChirpRequest(req.Body, &reqChirp)
	w.Header().Add(contentType, plainTextContentType)
	if errorString != "" || !isValid {
		returnErrorResponse(w, errorString)
		return
	}

	chirp := Chirp{
		Id:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      reqChirp.Body,
		UserId:    jwtId,
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
	authIdString := req.URL.Query().Get("author_id")
	sortString := req.URL.Query().Get("sort")
	sortAsc := sortString == "" || sortString == "asc" || sortString != "desc"
	if sortString == "desc" {
		sortAsc = false
	}
	if authIdString != "" {
		authorId, err := uuid.Parse(authIdString)
		if err != nil || authIdString == "" || authorId == uuid.Nil {
			returnChirpsResponse(w, []database.Chirp{})
			return
		}
		if sortAsc {
			cfg.getChirpsByAuthorIdAsc(w, req, authorId)
			return
		}
		cfg.getChirpsByAuthorIdDesc(w, req, authorId)
		return
	}
	if sortAsc {
		cfg.getAllChirpsAsc(w, req)
		return
	}
	cfg.getAllChirpsDesc(w, req)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("id")
	if idString == "" {
		returnNotFound(w)
	}
	chirpId, err := uuid.Parse(idString)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	dbChirp, err := cfg.dbQueries.GetChirpById(req.Context(), chirpId)
	if err != nil {
		returnNotFound(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add(contentType, plainTextContentType)
	err = json.NewEncoder(w).Encode(dbChirpToResponse(dbChirp))
	if err != nil {
		returnErrorResponse(w, standardError)
	}
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

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnUnauthorized(w)
		return
	}

	jwtId, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		if err.Error() == "invalid token" || err.Error() == "subject is empty" {
			returnBadRequest(w)
		} else {
			returnUnauthorized(w)
		}
		return
	}

	idString := req.PathValue("id")
	if idString == "" {
		returnNotFound(w)
	}
	chirpId, err := uuid.Parse(idString)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	dbChirp, err := cfg.dbQueries.GetChirpById(req.Context(), chirpId)
	if err != nil {
		returnNotFound(w)
		return
	}

	if dbChirp.UserID != jwtId {
		returnForbidden(w)
		return
	}
	_, err = cfg.dbQueries.DeleteChirpById(req.Context(), dbChirp.ID)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (cfg *apiConfig) getChirpsByAuthorIdAsc(w http.ResponseWriter, req *http.Request, authorId uuid.UUID) {
	dbChirps, err := cfg.dbQueries.GetChirpByAuthorIdAsc(req.Context(), authorId)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	returnChirpsResponse(w, dbChirps)
}

func (cfg *apiConfig) getChirpsByAuthorIdDesc(w http.ResponseWriter, req *http.Request, authorId uuid.UUID) {
	dbChirps, err := cfg.dbQueries.GetChirpByAuthorIdDesc(req.Context(), authorId)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	returnChirpsResponse(w, dbChirps)
}

func (cfg *apiConfig) getAllChirpsAsc(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.dbQueries.GetAllChirpsAsc(req.Context())
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	returnChirpsResponse(w, dbChirps)
}

func (cfg *apiConfig) getAllChirpsDesc(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.dbQueries.GetAllChirpsDesc(req.Context())
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	returnChirpsResponse(w, dbChirps)
}

func returnChirpsResponse(w http.ResponseWriter, chirps []database.Chirp) {
	responseChirps := []Chirp{}
	if len(chirps) > 0 {
		responseChirps = dbChirpsToResponse(chirps)
	}
	w.Header().Set(contentType, plainTextContentType)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(responseChirps)
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}
	w.Header().Add(contentType, plainTextContentType)
}
