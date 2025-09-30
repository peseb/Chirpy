package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/peseb/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getMetrics() int32 {
	return cfg.fileserverHits.Load()
}

func (cfg *apiConfig) resetMetrics() {
	cfg.fileserverHits.Store(0)
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	dbQueries := database.New(db)
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	// File/App routes
	fileHandler := http.FileServer(http.Dir("."))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileHandler)))

	// Api routes
	serveMux.HandleFunc("GET /api/healthz", handlerHealth)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	// Admin routes
	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(http.StatusOK)

	template := `<html>
					<body>
						<h1>Welcome, Chirpy Admin</h1>
						<p>Chirpy has been visited %d times!</p>
					</body>
				</html>`
	res := fmt.Sprintf(template, cfg.getMetrics())
	rw.Write([]byte(res))
}

func (cfg *apiConfig) handlerReset(rw http.ResponseWriter, req *http.Request) {
	cfg.resetMetrics()
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerHealth(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}

func handlerValidate(rw http.ResponseWriter, req *http.Request) {
	type validateDto struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	dto := validateDto{}
	err := decoder.Decode(&dto)
	if err != nil {
		respondWithError(rw, 400, "Unable to decode DTO", err)
		return
	}

	if len(dto.Body) > 140 {
		respondWithError(rw, 400, "Chirp is too long", nil)
		return
	}

	type validateResultDto struct {
		CleanedBody string `json:"cleaned_body"`
	}
	cleanedBody := cleanInput(dto.Body)
	respondWithJSON(rw, 200, validateResultDto{CleanedBody: cleanedBody})
}
