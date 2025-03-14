package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	hitCount := cfg.fileServerHits.Load()
	message := fmt.Sprintf(metricsHTML, hitCount)
	w.Write([]byte(message))
}

func (cfg *apiConfig) resetHitCountHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(int32(0))
	w.Write([]byte("Hit count successfully reset"))
}
