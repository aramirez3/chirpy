package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add(contentType, textHtmlContentType)

	usersCount, err := cfg.dbQueries.GetUsersCount(req.Context())
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	chirpsCount, err := cfg.dbQueries.GetChirpsCount(req.Context())
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(`<html>
	<head><style>body{background-color:#111;color:white;}</style></head>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
	<p>Total Users: %v</p>
	<p>Total Chirps: %v<p>
  </body>
</html>`, cfg.fileServerHits.Load(), usersCount, chirpsCount)
	w.Write([]byte(body))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits = atomic.Int32{}
	w.Header().Add(contentType, plainTextContentType)

	err := cfg.dbQueries.DeleteAllUsers(req.Context())
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	err = cfg.dbQueries.DeleteAllChirps(req.Context())
	if err != nil {
		returnErrorResponse(w, standardError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
