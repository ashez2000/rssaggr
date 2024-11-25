package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ashez2000/rssaggr/internal/auth"
	"github.com/ashez2000/rssaggr/internal/database"
	"github.com/google/uuid"
)

type Application struct {
	DB *database.Queries
}

func (app *Application) hello(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, "Hello, world!")
}

// createUser handler
// path: POST /users
func (app *Application) createUser(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Params{}
	err := decoder.Decode(&params)
	if err != nil {
		writeJSON(w, 400, "Error parsing body")
		return
	}

	user, err := app.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Username:  params.Username,
		ApiKey:    uuid.New().String(),
		CreatedAt: time.Now().UTC(),
	})

	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error creating user")
		return
	}

	writeJSON(w, 201, user)
}

// getUser handler
// path: GET /users
func (app *Application) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	writeJSON(w, 200, user)
}

// createFeed handler
// path: POST /feeds
func (app *Application) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type Params struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Params{}
	err := decoder.Decode(&params)
	if err != nil {
		writeJSON(w, 400, "Error parsing body")
		return
	}

	feed, err := app.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.URL,
		CreatedAt: time.Now().UTC(),
		UserID:    user.ID,
	})

	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error creating feed")
		return
	}

	writeJSON(w, 201, feed)
}

// getFeeds handler
// path: GET /feeds
func (app *Application) getFeeds(w http.ResponseWriter, r *http.Request) {
	feed, err := app.DB.GetFeeds(r.Context())
	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error getting feeds")
		return
	}

	writeJSON(w, 201, feed)
}

// createFeedFollow handler
// path: POST /feed_follows
func (app *Application) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type Params struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Params{}
	err := decoder.Decode(&params)
	if err != nil {
		writeJSON(w, 400, "Error parsing body")
		return
	}

	feed_follow, err := app.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
	})

	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error creating feed follow")
		return
	}

	writeJSON(w, 201, feed_follow)
}

// getFeedFollows handler
//
// path: GET /feed-follows
func (app *Application) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follows, err := app.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error getting feed_follows")
		return
	}

	writeJSON(w, 200, feed_follows)
}

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

// TODO: reafactor middleware
func (app *Application) authMiddleware(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r)
		if err != nil {
			writeJSON(w, 403, err)
			return
		}

		user, err := app.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			writeJSON(w, 403, err)
			return
		}

		handler(w, r, user)
	}
}
