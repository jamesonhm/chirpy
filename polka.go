package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jamesonhm/chirpy/internal/auth"
)

func (cfg *apiConfig) upgradeHandler(w http.ResponseWriter, r *http.Request) {
	type event struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		errorResponse(w, r, http.StatusUnauthorized, "unable to get api keky from header", err)
		return
	}
	if key != cfg.polkaKey {
		errorResponse(w, r, http.StatusUnauthorized, "wrong api key", nil)
		return
	}

	e, err := decode[event](r)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "invalid input", err)
		return
	}

	if e.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = cfg.db.UpgradeChirpyRed(r.Context(), e.Data.UserID)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, "unable to upgrade user", err)
		return
	}

	w.WriteHeader(204)
	return
}
