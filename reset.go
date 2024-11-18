package main

import (
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.serverHits.Swap(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
