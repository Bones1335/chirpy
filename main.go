package main

import (
	"encoding/json"
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

	mx.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		body := "OK"

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	})

	mx.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mx.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mx.HandleFunc("POST /api/validate_chirp", handlerJSON)

	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())))
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

func handlerJSON(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		jsonError := errorResponse{
			Error: "Something went wrong",
		}

		dat, _ := json.Marshal(jsonError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		jsonError := errorResponse{
			Error: "Chirp is too long",
		}

		dat, _ := json.Marshal(jsonError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	type successResponse struct {
    	Valid bool `json:"valid"`
	}
	
	jsonSuccess := successResponse{
		Valid: true,
	}
	dat, err := json.Marshal(jsonSuccess)

	if err != nil {
		jsonError := errorResponse{
			Error: "not valid",
		}
		dat, _ := json.Marshal(jsonError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}