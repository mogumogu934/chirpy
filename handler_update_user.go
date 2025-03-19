package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mogumogu934/chirpy/internal/auth"
	"github.com/mogumogu934/chirpy/internal/database"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding body: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Constraint, "email") {
				respondWithError(w, http.StatusBadRequest, "Email is already in use")
				return
			}
		}
		log.Printf("error updating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	dbUser, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("error getting user by email: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	type updateUserResp struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	user := updateUserResp{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, user)
}
