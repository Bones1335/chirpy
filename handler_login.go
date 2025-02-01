package main

import (
	"encoding/json"
	"net/http"

	"github.com/Bones1335/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	respondWithJSON(w, http.StatusOK, User{
		ID:        login.ID,
		CreatedAt: login.CreatedAt,
		UpdatedAt: login.UpdatedAt,
		Email:     login.Email,
	})
}
