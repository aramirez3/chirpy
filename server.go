package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	Handler *http.ServeMux
	Addr    string
}

const (
	serverRootPath = http.Dir("./src")
)

func createServer(port string) *Server {
	return &Server{http.NewServeMux(), ":" + port}
}

func startServer() {
	server := createServer("8080")

	server.Handler.Handle("/app/", http.StripPrefix("/app/", http.FileServer(serverRootPath)))
	server.Handler.HandleFunc("/healthz", handleReadiness)

	fmt.Printf("🐣 Chirping on http://localhost%s\n", server.Addr)
	err := http.ListenAndServe(server.Addr, server.Handler)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
