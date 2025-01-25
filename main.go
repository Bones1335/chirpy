package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mx := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mx,
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mx.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mx.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		body := "OK"

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	})

	mx.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mx.HandleFunc("POST /reset", apiCfg.handlerReset)

	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
