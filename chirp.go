package main

import (
	"net/http"
	"strings"
)

func deprofane(input string) string {
	profane := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(input, " ")
	for i, word := range words {
		if _, ok := profane[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func chirpValidHandler(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type retVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	c, err := decode[chirp](req)
	if err != nil {
		msg := "error decoding chirp"
		errorResponse(w, req, 500, msg, err)
		return
	}
	if len(c.Body) > 140 {
		msg := "Chirp too long"
		errorResponse(w, req, 400, msg, nil)
		return
	}

	respBody := retVal{
		CleanedBody: deprofane(c.Body),
	}
	encodeJsonResp(w, req, 200, respBody)
}
