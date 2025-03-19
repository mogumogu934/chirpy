package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mogumogu934/chirpy/internal/auth"
	"github.com/mogumogu934/chirpy/internal/database"
)

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
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting token string: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			respondWithError(w, http.StatusUnauthorized, "Token has expired")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			respondWithError(w, http.StatusUnauthorized, "Invalid token format - access token required")
		default:
			respondWithError(w, http.StatusUnauthorized, "Authentication failed")
		}
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding body: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   filterMsg(params.Body),
		UserID: userID,
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("error creating chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp")
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
