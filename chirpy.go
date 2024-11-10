package main

import (
	"fmt"
	"net/http"
)

type requestHandler struct{}

func (requestHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

type Server struct {
	Server  http.Server
	Handler http.ServeMux
}

func startServer() {

	server := Server{}
	mux := http.NewServeMux()
	mux.Handle("/", requestHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "🐣 Chirpy open on http://loaclhost:8080")
	})
}
