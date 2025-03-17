package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) resetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	if err := cfg.db.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("error deleting all users: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error deleting all users")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "All users successfully deleted"})
}
