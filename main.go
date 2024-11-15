package main

import (
	"log"
	"net/http"
)

func main() {
	const port = ":8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
