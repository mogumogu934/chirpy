package main

import (
	"log"
	"net/http"
)

const metricsHTML = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func main() {
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitCountHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHitCountHandler)

	port := "8080"
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
