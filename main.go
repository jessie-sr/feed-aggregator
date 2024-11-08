package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jessie-sr/rss-aggregator/internal/db"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // Include this even through I'm not calling it directly
)

type apiConfig struct {
	DB *db.Queries // Hold connection to the database
}

func main() {
	// Load variables in .env
	godotenv.Load()

	// Get port number
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in the environment")
	}

	// Get db url
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Database url is not found in the environment")
	}

	// Connect to db
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	apiCig := apiConfig{
		DB: db.New(conn),
	}

	// Create a new router
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Connect handleReadiness() to the different paths
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness)
	v1Router.Get("/error", handleError)
	v1Router.Post("/users", apiCig.handleCreateUser)

	// Nesting v1Router under the main router
	router.Mount("/v1", v1Router) // "/healthz" -> "/v1/healthz"

	// Create a new server with the router and port number
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
