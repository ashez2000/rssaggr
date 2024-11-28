package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ashez2000/rssaggr/internal/auth"
	"github.com/ashez2000/rssaggr/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type application struct {
	store *database.Queries
}

func (app *application) hello(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, "Hello, world!")
}

// createUser handler
// path: POST /users
func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := app.store.CreateUser(r.Context(), database.CreateUserParams{
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
func (app *application) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	writeJSON(w, 200, user)
}

// createFeed handler
// path: POST /feeds
func (app *application) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
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

	feed, err := app.store.CreateFeed(r.Context(), database.CreateFeedParams{
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
func (app *application) getFeeds(w http.ResponseWriter, r *http.Request) {
	feed, err := app.store.GetFeeds(r.Context())
	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error getting feeds")
		return
	}

	writeJSON(w, 201, feed)
}

// createFeedFollow handler
// path: POST /feed_follows
func (app *application) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
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

	feed_follow, err := app.store.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
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
func (app *application) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follows, err := app.store.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error getting feed_follows")
		return
	}

	writeJSON(w, 200, feed_follows)
}

// deleteFeedFollow handler
//
// path: DELETE /feed-follows
func (app *application) deleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedID, err := uuid.Parse(chi.URLParam(r, "feedID"))
	if err != nil {
		writeJSON(w, 500, "Error parsing feed follow id")
		return
	}

	err = app.store.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feedID,
		UserID: user.ID,
	})
	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error getting feed_follows")
		return
	}

	writeJSON(w, 200, "feed follow removed")
}

// getPostsForUser handler
//
// path: GET /posts
func (app *application) getPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := app.store.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "error getting posts")
		return
	}

	writeJSON(w, 200, posts)
}

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

// TODO: reafactor middleware
func (app *application) authMiddleware(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r)
		if err != nil {
			writeJSON(w, 403, err)
			return
		}

		user, err := app.store.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			writeJSON(w, 403, err)
			return
		}

		handler(w, r, user)
	}
}
