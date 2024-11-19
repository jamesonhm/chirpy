package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type userEmail struct {
		Email string `json:"email"`
	}
	type resp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	params, err := decode[userEmail](r)
	if err != nil {
		errorResponse(w, r, 500, "error decoding user", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		errorResponse(w, r, 500, "error creating user", err)
	}
	userResp := resp(user)
	err = encodeJsonResp(w, r, 201, userResp)

}
