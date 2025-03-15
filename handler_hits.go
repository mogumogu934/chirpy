package main

import (
	"fmt"
	"net/http"
)

const metricsHTML = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden: metrics only available in dev environment")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	hitCount := cfg.fileServerHits.Load()
	message := fmt.Sprintf(metricsHTML, hitCount)
	w.Write([]byte(message))
}

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden: reset metrics only available in dev environment")
		return
	}

	cfg.fileServerHits.Store(int32(0))
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hit count successfully reset"})
}
