package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/mogumogu934/chirpy/internal/auth"
	"github.com/mogumogu934/chirpy/internal/database"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	tokenOldString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting token string: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	dbToken, err := cfg.db.GetRefreshToken(r.Context(), tokenOldString)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		log.Printf("error getting token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if time.Now().After(dbToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Expired token")
		return
	}

	tokenString, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error making refresh token")
		respondWithError(w, http.StatusInternalServerError, "Error refreshing token")
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  tokenString,
		UserID: dbToken.UserID,
	})

	if err != nil {
		log.Printf("error making refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error refreshing token")
		return
	}

	accessToken, err := auth.MakeJWT(
		dbToken.UserID,
		cfg.jwtSecret,
		time.Hour,
	)

	if err != nil {
		log.Printf("error making JWT: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error refreshing token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": accessToken,
	})
}
