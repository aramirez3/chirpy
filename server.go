package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Handler *http.ServeMux
	Addr    string
	Config  apiConfig
}

const (
	serverRootPath    = http.Dir("./src")
	contentType       = "Content-Type"
	headerContentType = "text/plain; charset=utf-8"
)

func createServer(port string) *Server {
	return &Server{http.NewServeMux(), ":" + port, apiConfig{fileServerHits: atomic.Int32{}}}
}

func startServer() {
	server := createServer("8080")

	server.Handler.Handle("/app/", http.StripPrefix("/app/", server.Config.middlewareMetricsInc(http.FileServer(serverRootPath))))
	server.Handler.HandleFunc("/healthz", handleReadiness)
	server.Handler.HandleFunc("/metrics", server.Config.handlerMetrics)

	fmt.Printf("🐣 Chirping on http://localhost%s\n", server.Addr)
	err := http.ListenAndServe(server.Addr, server.Handler)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add(contentType, headerContentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add(contentType, headerContentType)
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileServerHits.Load()
	w.Write([]byte(strconv.Itoa(int(hits))))
}
