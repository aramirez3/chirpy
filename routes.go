package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
)

const (
	contentType       = "Content-Type"
	headerContentType = "text/plain; charset=utf-8"
)

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add(contentType, headerContentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add(contentType, headerContentType)
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf("Hits: %v\n", strconv.Itoa(int(cfg.fileServerHits.Load())))
	w.Write([]byte(hits))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits = atomic.Int32{}
	w.Header().Add(contentType, headerContentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
