package main

import (
	"encoding/json"
	"log"
	"net/http"
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

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		returnErrorResponse(w)
		return
	}

	w.Header().Add(contentType, plainTextContentType)

	if len(chirp.Body) <= 140 {
		w.WriteHeader(http.StatusOK)
		validResponse := ValidResponse{
			Valid: true,
		}
		respBody, _ := encodeJson(validResponse)
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

func encodeJson(body any) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		log.Printf("error marshaling json: %s\n", err)
	}
	return data, nil
}

func returnErrorResponse(w http.ResponseWriter) {
	w.Header().Add(contentType, plainTextContentType)
	w.WriteHeader(http.StatusBadRequest)
	errorResponse := ErrorResponse{
		Error: "Something went wrong",
	}
	respBody, _ := encodeJson(errorResponse)
	w.Write(respBody)

}
