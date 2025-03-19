package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mogumogu934/chirpy/internal/auth"
	"github.com/mogumogu934/chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		ExpiresIn int    `json:"expires_in_seconds,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding body: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error logging in")
		return
	}

	if params.Password == "" {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if params.ExpiresIn == 0 {
		params.ExpiresIn = 3600 // Default to 1 hour
	} else if params.ExpiresIn > 3600 {
		params.ExpiresIn = 3600 // Limit to 1 hour
	}

	token, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		time.Duration(params.ExpiresIn),
	)

	if err != nil {
		log.Printf("error making JWT: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error logging in")
		return
	}

	refTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error making refresh token string: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error logging in")
		return
	}

	refToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refTokenString,
		UserID: dbUser.ID,
	})

	if err != nil {
		log.Printf("error making refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error logging in")
		return
	}

	type loginUserResp struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	user := loginUserResp{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refToken.Token,
	}

	respondWithJSON(w, http.StatusOK, user)
}
