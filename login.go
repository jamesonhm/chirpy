package main

import (
	"net/http"
	"time"

	"github.com/jamesonhm/chirpy/internal/auth"
)

type response struct {
	userResp
	Token string `json:"token"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type creds struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600
	}
	expireTime := time.Second * time.Duration(params.ExpiresInSeconds)
	token, err := auth.MakeJWT(userDb.ID, cfg.tokenSecret, expireTime)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "error creating jwt", err)
	}

	encodeJsonResp(w, r, http.StatusOK, response{
		userResp: userResp{
			ID:        userDb.ID,
			CreatedAt: userDb.CreatedAt,
			UpdatedAt: userDb.UpdatedAt,
			Email:     userDb.Email,
		},
		Token: token,
	})
}
