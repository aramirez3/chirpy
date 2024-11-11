package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Server struct {
	Handler *http.ServeMux
	Addr    string
	Config  apiConfig
}

const (
	serverRootPath       = http.Dir("./src")
	adminPath            = http.Dir("./admin")
	contentType          = "Content-Type"
	plainTextContentType = "text/plain; charset=utf-8"
	textHtmlContentType  = "text/html; charset=utf-8"
)

func createServer(port string) *Server {
	return &Server{http.NewServeMux(), ":" + port, apiConfig{fileServerHits: atomic.Int32{}}}
}

func (s *Server) startServer() {
	s.Handler.Handle("/app/", http.StripPrefix("/app/", s.Config.middlewareMetricsInc(http.FileServer(serverRootPath))))
	s.Handler.HandleFunc("GET /api/healthz", handleReadiness)
	s.Handler.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	s.Handler.HandleFunc("GET /admin/metrics", s.Config.handlerMetrics)
	s.Handler.HandleFunc("POST /admin/reset", s.Config.handleReset)

	fmt.Printf("🐣 Chirping on http://localhost%s\n", s.Addr)
	err := http.ListenAndServe(s.Addr, s.Handler)

	if err != nil {
		fmt.Println(err)
		return
	}
}
