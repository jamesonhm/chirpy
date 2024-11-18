package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/jamesonhm/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	serverHits atomic.Int32
	queries    *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	const port = ":8080"
	const filepathRoot = "."

	cfg := &apiConfig{
		serverHits: atomic.Int32{},
		queries:    database.New(db),
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("POST /api/validate_chirp", chirpValidHandler)

	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
