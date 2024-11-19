package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("error decoding json: %w", err)
	}
	return v, nil
}

func encodeJsonResp(w http.ResponseWriter, r *http.Request, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(v)
	if err != nil {
		msg := fmt.Sprintf("Error marshaling JSON: %s", err)
		log.Printf(msg)
		w.WriteHeader(500)
		return fmt.Errorf(msg)
	}
	w.WriteHeader(status)
	w.Write(dat)
	return nil
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if status > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	_ = encodeJsonResp(w, r, status, errorResponse{
		Error: msg,
	})
}
