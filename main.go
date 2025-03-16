package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mogumogu934/chirpy/internal/database"
)

type apiConfig struct {
	db             *database.Queries
	platform       string
	fileServerHits atomic.Int32
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading environment")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("error opening connection to database")
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db:       dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIDHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.getMetricsHandler)
	mux.HandleFunc("POST /admin/reset-metrics", apiCfg.resetMetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	port := "8080"
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
