package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mogumogu934/chirpy/internal/auth"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
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

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("error decoding chirpID string into UUID: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting chirp")
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Chirp with ID %v does not exist", chirpID))
			return
		}
		log.Printf("error getting chirp by ID: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting chirp")
		return
	}

	if userID != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, "Chirp can only be deleted by author")
		return
	}

	if err = cfg.db.DeleteChirp(r.Context(), dbChirp.ID); err != nil {
		log.Printf("error deleting chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
