package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func validateLengthHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	type validResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		writeJSONResponse(w, 500, errorResponse{Error: "Something went wrong"})
		return
	}

	if len(params.Body) > 140 {
		writeJSONResponse(w, 400, errorResponse{Error: "Chirp is too long"})
		return
	}

	writeJSONResponse(w, http.StatusOK, validResponse{Valid: true})
}
