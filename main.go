package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/kevin120202/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	db 				*database.Queries
	platform		string
	secretKey		string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platformCfg := os.Getenv("PLATFORM")
	if platformCfg == "" {
		log.Fatal("PLATFORM must be set")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db: 			dbQueries,
		platform: 		platformCfg,
		secretKey: 		secretKey,
	}

	mux := http.NewServeMux()
	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)
	
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.reset)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiConfig.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiConfig.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerChirpsGetOne)
	mux.HandleFunc("POST /api/users", apiConfig.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiConfig.handlerUsersLogin)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port:: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

