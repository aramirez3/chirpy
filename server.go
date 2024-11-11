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

func startServer() {
	server := createServer("8080")

	server.Handler.Handle("/app/", http.StripPrefix("/app/", server.Config.middlewareMetricsInc(http.FileServer(serverRootPath))))
	server.Handler.HandleFunc("GET /api/healthz", handleReadiness)
	server.Handler.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	server.Handler.HandleFunc("GET /admin/metrics", server.Config.handlerMetrics)
	server.Handler.HandleFunc("POST /admin/reset", server.Config.handleReset)

	fmt.Printf("🐣 Chirping on http://localhost%s\n", server.Addr)
	err := http.ListenAndServe(server.Addr, server.Handler)

	if err != nil {
		fmt.Println(err)
		return
	}
}
