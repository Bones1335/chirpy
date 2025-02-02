package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Bones1335/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string         `json:"email"`
		Password         string         `json:"password"`
		ExpiresInSeconds *time.Duration `json:"expires_in_seconds"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	login, err := cfg.db.Login(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Email not found", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, login.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	const DefaultExpiry = 1 * time.Hour
	var expiry time.Duration

	if params.ExpiresInSeconds == nil {
		expiry = DefaultExpiry
	} else if *params.ExpiresInSeconds > DefaultExpiry {
		expiry = DefaultExpiry
	} else {
		expiry = *params.ExpiresInSeconds
	}

	token, err := auth.MakeJWT(login.ID, cfg.jwtSecret, expiry)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        login.ID,
		CreatedAt: login.CreatedAt,
		UpdatedAt: login.UpdatedAt,
		Email:     login.Email,
		Token:     token,
	})
}
