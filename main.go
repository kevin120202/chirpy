package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Add(1)

			next.ServeHTTP(w, r)
		},
	)
}

func (cfg *apiConfig) writeNumOfRequest(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	msg := fmt.Sprintf("Hits: %d", count)
	w.Write([]byte(msg))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	mux.HandleFunc("GET /healthz", readiness)
	mux.HandleFunc("GET /metrics", apiCfg.writeNumOfRequest)
	mux.HandleFunc("POST /reset", apiCfg.reset)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}