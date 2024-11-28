package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ashez2000/rssaggr/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

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
