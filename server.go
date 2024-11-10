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

func startServer() {
	server := &Server{http.NewServeMux(), ":8080"}

	server.Handler.Handle("/", http.FileServer(serverRootPath))

	fmt.Printf("🐣 Chirping on http://localhost%s\n", server.Addr)
	err := http.ListenAndServe(server.Addr, server.Handler)

	if err != nil {
		fmt.Println(err)
		return
	}
}
