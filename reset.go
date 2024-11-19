package main

import (
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		errorResponse(w, r, 403, "dev only function", nil)
		return
	}
	cfg.serverHits.Swap(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		errorResponse(w, r, 500, "error deleting users", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state"))
}
