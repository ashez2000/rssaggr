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

func (app *Application) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	writeJSON(w, 200, user)
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
