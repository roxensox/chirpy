package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/roxensox/chirpy/internal/chirpyserver"
	"github.com/roxensox/chirpy/internal/database"
	"net/http"
	"os"
)

func main() {
	// Loads in the .env file
	godotenv.Load()

	// Gets the database URL from .env
	dbURL := os.Getenv("DB_URL")

	// Opens the database and handles errors
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Unable to open database")
		os.Exit(1)
	}

	// Gets a query engine for the database and adds it to the config object
	dbQueries := database.New(db)
	config := chirpyserver.ApiConfig{
		DBConn: dbQueries,
		Secret: os.Getenv("SECRET"),
	}

	// Starts a new server mux
	sMux := http.NewServeMux()

	// Makes a root handler for the file server
	handler := http.FileServer(http.Dir("."))

	// Sets up a server object
	server := http.Server{
		Handler: sMux,
		Addr:    ":8080",
	}

	// Binds middleware-produced handler to app directory
	sMux.Handle("/app/", config.MiddlewareMetricsInc(handler))

	// Binds functions to POST handlers
	sMux.HandleFunc("POST /admin/reset", config.Reset)
	sMux.HandleFunc("POST /api/validate_chirp", chirpyserver.ValidateChirp)
	sMux.HandleFunc("POST /api/users", config.POSTUsers)
	sMux.HandleFunc("POST /api/chirps", config.POSTChirps)
	sMux.HandleFunc("POST /api/login", config.POSTLogin)
	sMux.HandleFunc("POST /api/refresh", config.POSTRefresh)
	sMux.HandleFunc("POST /api/revoke", config.POSTRevoke)

	// Binds functions to GET handlers
	sMux.HandleFunc("GET /api/healthz", chirpyserver.Healthz)
	sMux.HandleFunc("GET /admin/metrics", config.FServerHits)
	sMux.HandleFunc("GET /api/chirps", config.GETChirps)
	sMux.HandleFunc("GET /api/chirps/{chirpID}", config.GETChirpByID)

	// Runs the server
	server.ListenAndServe()
}
