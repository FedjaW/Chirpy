package main

import (
    "github.com/FedjaW/Chirpy/internal/database"
    "os"
    "database/sql"
    "github.com/joho/godotenv"
	"log"
	"net/http"
	"sync/atomic"
    _ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
    db *database.Queries
    platform string
}

func main() {
    godotenv.Load()
    platform := os.Getenv("PLATFORM")
    if platform == "" {
        log.Fatalf("PLATFORM must be set")
    }
    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatalf("DB_URL must be set")
    }
    dbConn, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Error opening database %s", err)
    }
    dbQueries := database.New(dbConn)

	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
        fileserverHits: atomic.Int32{},
        db: dbQueries,
		platform: platform,
    }

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
    mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
    mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
