package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

const (
	contentType          = "Content-Type"
	plainTextContentType = "text/plain; charset=utf-8"
	textHtmlContentType  = "text/html; charset=utf-8"
)

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add(contentType, plainTextContentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add(contentType, textHtmlContentType)
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileServerHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits = atomic.Int32{}
	w.Header().Add(contentType, plainTextContentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
