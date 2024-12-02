package main

import (
	"net/http"
	"time"

	"github.com/jamesonhm/chirpy/internal/auth"
	"github.com/jamesonhm/chirpy/internal/database"
)

const (
	jwt_expire_seconds = 3600
	expireTime         = time.Second * time.Duration(jwt_expire_seconds)
)

type response struct {
	userResp
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

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

	token, err := auth.MakeJWT(userDb.ID, cfg.tokenSecret, expireTime)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "error creating jwt", err)
		return
	}
	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}
	_, err = cfg.db.CreateRefresh(r.Context(), database.CreateRefreshParams{
		Token:     refresh_token,
		UserID:    userDb.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	encodeJsonResp(w, r, http.StatusOK, response{
		userResp: userResp{
			ID:          userDb.ID,
			CreatedAt:   userDb.CreatedAt,
			UpdatedAt:   userDb.UpdatedAt,
			Email:       userDb.Email,
			IsChirpyRed: userDb.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refresh_token,
	})
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, r, http.StatusUnauthorized, "unable to get jwt from header", err)
		return
	}
	tokenDb, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		errorResponse(w, r, 401, "not logged in", err)
		return
	}

	if time.Now().After(tokenDb.ExpiresAt) || (tokenDb.RevokedAt.Valid && time.Now().After(tokenDb.RevokedAt.Time)) {
		errorResponse(w, r, 401, "not logged in", err)
		return
	}

	jwt, err := auth.MakeJWT(tokenDb.UserID, cfg.tokenSecret, expireTime)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "error creating jwt", err)
		return
	}

	encodeJsonResp(w, r, http.StatusOK, response{
		Token: jwt,
	})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, r, http.StatusUnauthorized, "unable to get jwt from header", err)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		errorResponse(w, r, 401, "not logged in", err)
		return
	}

	w.WriteHeader(204)

}
