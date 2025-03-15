package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, validResponse{CleanedBody: filterMsg(params.Body)})
}
