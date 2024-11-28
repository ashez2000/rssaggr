package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ashez2000/rssaggr/internal/database"
	"github.com/google/uuid"
)

func (app *application) getFeeds(w http.ResponseWriter, r *http.Request) {
	feed, err := app.store.GetFeeds(r.Context())
	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error getting feeds")
		return
	}

	writeJSON(w, 201, feed)
}

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
