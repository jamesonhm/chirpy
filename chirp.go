package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jamesonhm/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	p, err := decode[params](r)
	if err != nil {
		msg := "error decoding chirp"
		errorResponse(w, r, 500, msg, err)
		return
	}
	cleanedBody, err := validateChirp(p.Body)
	if err != nil {
		errorResponse(w, r, 400, "error validating chirp", err)
	}

	createdChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: p.UserID,
	})
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "error creating chirp", err)
		return
	}
	respBody := Chirp(createdChirp)
	encodeJsonResp(w, r, 201, respBody)

}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("Chirp too long")
	}

	profane := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := profane[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " "), nil
}
