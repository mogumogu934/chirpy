package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", apiCfg.hitCountHandler)
	mux.HandleFunc("/reset", apiCfg.resetHitCountHandler)

	port := "8080"
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitCountHandler(w http.ResponseWriter, r *http.Request) {
	hitCount := cfg.fileServerHits.Load()
	message := fmt.Sprintf("Hits: %d", hitCount)
	w.Write([]byte(message))
}

func (cfg *apiConfig) resetHitCountHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(int32(0))
	w.Write([]byte("Hit count successfully reset"))
}
