package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jessie-sr/rss-aggregator/internal/config"
	"github.com/jessie-sr/rss-aggregator/internal/db"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // Include this even through I'm not calling it directly
)

type apiConfig struct {
	DB *db.Queries // Hold connection to the database
}

func main() {
	/*
		To test urlToFeed:
		feeds, err := urlToFeed("https://rss.nytimes.com/services/xml/rss/nyt/PersonalTech.xml")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(feeds)
	*/

	// Load variables in .env
	godotenv.Load()

	// Get port number
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in the environment")
	}

	filePath, err := config.GetFilePath()
	if err != nil {
		log.Fatal("Can't find the file", err)
	}

	// Read the Config struct located at ~/.gatorconfig.json
	cfg, err := config.Read(filePath)
	if err != nil {
		log.Fatal("Can't read the file", err)
	}

	// Store config in a state
	state := State{
		Ptr: &cfg,
	}

	// Initialize cmds and register handler functions
	cmds := Commands{}
	cmds.register("login", handlerLogin) // Register the "login" command

	// Validate CLI input
	if len(os.Args) < 2 {
		log.Fatal("Malformed command: no command provided")
	}

	// Parse CLI input to create a Command
	cmd := Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	// Execute the command
	err = cmds.run(&state, cmd)
	if err != nil {
		log.Fatalf("Error executing command '%s': %v", cmd.Name, err)
	}

	// Connect to db
	dbURL := cfg.DBUrl
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	database := db.New(conn)
	apiCig := apiConfig{
		DB: database,
	}

	go startScraping(database, 10, time.Minute)

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

	// Connect handlerReadiness() to the different paths
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)

	v1Router.Post("/users", apiCig.handlerCreateUser)
	v1Router.Get("/users", apiCig.middlewareAuth(apiCig.handlerGetUser)) // Use middlewareAuth to convert handlerGetUser to regular http.HandlerrFunc

	v1Router.Post("/feeds", apiCig.middlewareAuth(apiCig.handlerCreateFeed))
	v1Router.Get("/feeds", apiCig.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCig.middlewareAuth(apiCig.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiCig.middlewareAuth(apiCig.handlerGetFollowedFeeds))
	v1Router.Delete("/feed_follows/{feed_follows_id}", apiCig.middlewareAuth(apiCig.handlerUnfollowFeed))

	v1Router.Get("/posts", apiCig.middlewareAuth(apiCig.handlerGetPostsForUser))

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
