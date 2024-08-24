package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/vickon16/first_backend_tutorial/internal/database"

	_ "github.com/lib/pq" // include postgres without calling it
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB environment variable not set")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	db := database.New(conn)
	apiConfig := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	router.Mount("/v1", v1Router) // versioning the routers
	v1Router.Get("/ready", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", apiConfig.handlerCreateUser)
	v1Router.Get("/user", apiConfig.middleWareAuth(apiConfig.handlerGetUserByAPIKey))

	// Feed
	v1Router.Get("/feeds", apiConfig.handlerGetFeeds)
	v1Router.Post("/feeds", apiConfig.middleWareAuth(apiConfig.handlerCreateFeed))

	// Feed Follows
	v1Router.Post("/feed-follows", apiConfig.middleWareAuth(apiConfig.handlerCreateFeedFollows))
	v1Router.Get("/feed-follows", apiConfig.middleWareAuth(apiConfig.handlerGetFeedFollows))
	v1Router.Delete("/feed-follows/{feedFollowId}", apiConfig.middleWareAuth(apiConfig.handlerDeleteFeedFollow))

	// Post
	v1Router.Get("/posts", apiConfig.middleWareAuth(apiConfig.handlerGetPostsForUser))

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server Starting on port %s", portString)
	serverError := server.ListenAndServe()
	if serverError != nil {
		log.Fatal(serverError)
	}
}
