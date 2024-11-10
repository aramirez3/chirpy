package main

import (
	"fmt"
	"net/http"
)

type requestHandler struct{}

func (requestHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

type Server struct {
	Handler *http.ServeMux
	Addr    string
}

func startServer() {
	mux := http.NewServeMux()
	server := Server{mux, ":8080"}
	fmt.Printf("🐣 Chirping on http://localhost%s\n", server.Addr)
	err := http.ListenAndServe(server.Addr, server.Handler)
	if err != nil {
		fmt.Println(err)
		return
	}
}
