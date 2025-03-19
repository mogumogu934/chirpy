package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/mogumogu934/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting token string: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid or already revoked token")
			return
		}
		log.Printf("error revoking token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
