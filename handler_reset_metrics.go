package main

import "net/http"

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	cfg.fileServerHits.Store(int32(0))
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hit count successfully reset"})
}
