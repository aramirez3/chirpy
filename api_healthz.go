package main

import "net/http"

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add(contentType, plainTextContentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
