package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding body: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("User with ID %v does not exist", params.Data.UserID))
			return
		}
		log.Printf("error upgrading user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
