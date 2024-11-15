package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/aramirez3/chirpy/internal/database"
)

type Server struct {
	Handler *http.ServeMux
	Addr    string
	Config  apiConfig
}

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
	Secret         string
	PolkaKey       string
}

const (
	serverRootPath       = http.Dir("./src")
	adminPath            = http.Dir("./admin")
	contentType          = "Content-Type"
	plainTextContentType = "text/plain; charset=utf-8"
	textHtmlContentType  = "text/html; charset=utf-8"
	standardError        = "Something went wrong"
)

func createServer(port string) *Server {
	return &Server{http.NewServeMux(), ":" + port, apiConfig{fileServerHits: atomic.Int32{}}}
}

func (s *Server) startServer() {
	s.Handler.Handle("/app/", http.StripPrefix("/app/", s.Config.middlewareMetricsInc(http.FileServer(serverRootPath))))
	s.Handler.HandleFunc("GET /api/healthz", handleReadiness)
	s.Handler.HandleFunc("POST /api/chirps", s.Config.handleNewChirp)
	s.Handler.HandleFunc("GET /api/chirps", s.Config.handleGetChirps)
	s.Handler.HandleFunc("GET /api/chirps/{id}", s.Config.handleGetChirp)
	s.Handler.HandleFunc("DELETE /api/chirps/{id}", s.Config.handleDeleteChirp)
	s.Handler.HandleFunc("GET /admin/metrics", s.Config.handlerMetrics)
	s.Handler.HandleFunc("POST /admin/reset", s.Config.handleReset)
	s.Handler.HandleFunc("POST /api/users", s.Config.handleNewUser)
	s.Handler.HandleFunc("PUT /api/users", s.Config.handleUserUpdate)
	s.Handler.HandleFunc("POST /api/login", s.Config.handleLogin)
	s.Handler.HandleFunc("POST /api/refresh", s.Config.handleRefresh)
	s.Handler.HandleFunc("POST /api/revoke", s.Config.handleRevoke)
	s.Handler.HandleFunc("POST /api/polka/webhooks", s.Config.handlePolkaWebhooks)
	fmt.Printf("🐣 Chirping on http://localhost%s\n", s.Addr)
	err := http.ListenAndServe(s.Addr, s.Handler)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func encodeJson(body any) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		log.Printf("error marshaling json: %s\n", err)
	}
	return data, nil
}

func returnErrorResponse(w http.ResponseWriter, errorString string) {
	if errorString == "" {
		errorString = standardError
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set(contentType, plainTextContentType)
	respBody, _ := encodeJson(ErrorResponse{
		Error: errorString,
	})
	w.Write(respBody)
}

func returnNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Add(contentType, plainTextContentType)
	respBody, _ := encodeJson(ErrorResponse{
		Error: http.StatusText(http.StatusNotFound),
	})
	w.Write(respBody)
}

func returnUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Add(contentType, plainTextContentType)
	respBody, _ := encodeJson(ErrorResponse{
		Error: http.StatusText(http.StatusUnauthorized),
	})
	w.Write(respBody)
}

func returnBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add(contentType, plainTextContentType)
	respBody, _ := encodeJson(ErrorResponse{
		Error: http.StatusText(http.StatusBadRequest),
	})
	w.Write(respBody)
}

func returnForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add(contentType, plainTextContentType)
	respBody, _ := encodeJson(ErrorResponse{
		Error: http.StatusText(http.StatusForbidden),
	})
	w.Write(respBody)
}
