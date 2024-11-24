package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
		CreatedAt: time.Now().UTC(),
	})

	if err != nil {
		log.Println(err)
		writeJSON(w, 500, "Error creating user")
		return
	}

	writeJSON(w, 201, user)
}
