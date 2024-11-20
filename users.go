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
