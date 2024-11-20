package main

import (
	"net/http"

	"github.com/jamesonhm/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params, err := decode[creds](r)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "user input error", err)
		return
	}

	userDb, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, "error getting user by email", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, userDb.HashedPassword)
	if err != nil {
		errorResponse(w, r, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	encodeJsonResp(w, r, http.StatusOK, userResp{
		ID:        userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:     userDb.Email,
	})
}
