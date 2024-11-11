package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Chirp struct {
	Body string `json:"body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidResponse struct {
	Valid bool `json:"valid"`
}

type CleanChirp struct {
	CleanedBody string `json:"cleaned_body"`
}

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		returnErrorResponse(w)
		return
	}
	if chirp.Body == "" {
		returnErrorResponse(w)
		return
	}

	w.Header().Add(contentType, plainTextContentType)

	if len(chirp.Body) <= 140 {
		w.WriteHeader(http.StatusOK)
		respBody, _ := encodeJson(removeProfanity(chirp))
		w.Write(respBody)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ErrorResponse{
			Error: "Chirp is too long",
		}
		respBody, _ := encodeJson(errorResponse)
		w.Write(respBody)
	}

}

func removeProfanity(chirp Chirp) CleanChirp {
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	cleanChirp := CleanChirp{}
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
		cleanChirp.CleanedBody = strings.Join(words, " ")
	}
	return cleanChirp
}
