package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ashez2000/rssaggr/internal/database"
	"github.com/ashez2000/rssaggr/internal/rss"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	db := database.New(conn)
	app := Application{
		DB: db,
	}

	go rss.AggrRSSFeeds(db, 3, 10*time.Second)

	router := chi.NewRouter()

	// base middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

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
	router.Get("/posts", app.authMiddleware(app.getPostsForUser))

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
