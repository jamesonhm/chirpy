package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jamesonhm/chirpy/internal/auth"
	"github.com/jamesonhm/chirpy/internal/database"
)

type userResp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type userEntry struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params, err := decode[userEntry](r)
	if err != nil {
		errorResponse(w, r, 500, "error decoding user", err)
		return
	}
	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "error hashing pw", err)
	}

	userDb, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		errorResponse(w, r, 500, "error creating user", err)
	}
	user := userResp(userDb)
	err = encodeJsonResp(w, r, 201, user)

}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, r, http.StatusUnauthorized, "unable to get jwt from header", err)
		return
	}
	validID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		errorResponse(w, r, http.StatusUnauthorized, "unable to validate jwt", err)
		return
	}

	params, err := decode[creds](r)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "user input error", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "error hashing pw", err)
	}

	userDb, err := cfg.db.UpdateUserCreds(r.Context(), database.UpdateUserCredsParams{
		Email:          params.Email,
		HashedPassword: hashed,
		ID:             validID,
	})
	if err != nil {
		errorResponse(w, r, 500, "error updating user credentials", err)
	}

	err = encodeJsonResp(w, r, http.StatusOK, userResp(userDb))
}
