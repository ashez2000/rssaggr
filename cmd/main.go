package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ashez2000/rssaggr/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := loadEnv("PORT")
	databaseURL := loadEnv("DATABASE_URL")

	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	app := Application{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API routes
	router.Get("/", app.hello)
	router.Get("/users", app.authMiddleware(app.getUser))
	router.Post("/users", app.createUser)
	router.Post("/feeds", app.authMiddleware(app.createFeed))
	router.Get("/feeds", app.getFeeds)
	router.Post("/feed-follows", app.authMiddleware(app.createFeedFollow))
	router.Get("/feed-follows", app.authMiddleware(app.getFeedFollows))
	router.Delete("/feed-follows/{feedID}", app.authMiddleware(app.deleteFeedFollow))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Println("Listening on port", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func loadEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Sprintf("%v undefined", name))
	}

	return value
}
