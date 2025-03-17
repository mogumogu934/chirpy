package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mogumogu934/chirpy/internal/database"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	dBchirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting chirps: %v", err))
		return
	}

	chirps := []Chirp{}
	for _, chirp := range dBchirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error decoding chirpID string into UUID: %v", err))
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("chirp with chirpID %v does not exist", chirpID))
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting chirp by chirpID: %v", err))
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func filterMsg(msg string) string {
	blockedWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	split := strings.Split(msg, " ")
	for i, word := range split {
		lowercaseWord := strings.ToLower(word)
		if _, exists := blockedWords[lowercaseWord]; exists {
			split[i] = "****"
		}
	}

	return strings.Join(split, " ")
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error decoding body: %v", err))
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   filterMsg(params.Body),
		UserID: params.UserID,
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating chirp: %v", err))
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
