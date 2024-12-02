package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jamesonhm/chirpy/internal/auth"
	"github.com/jamesonhm/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, "Invalid chirp ID", nil)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, "chirp not found", err)
		return
	}

	resp := Chirp(chirp)
	encodeJsonResp(w, r, http.StatusOK, resp)

}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "error getting chirps from db", err)
		return
	}

	resp := []Chirp{}
	for _, chirp := range chirps {
		resp = append(resp, Chirp(chirp))
	}

	encodeJsonResp(w, r, http.StatusOK, resp)
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
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
		UserID: validID,
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

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
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

	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, "Invalid chirp ID", nil)
		return
	}
	chirpDb, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, "chirp not found", err)
		return
	}

	if chirpDb.UserID != validID {
		errorResponse(w, r, http.StatusForbidden, "you are not the author of this chirp", nil)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpUUID)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, "unable to delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
