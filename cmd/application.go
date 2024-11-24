package main

import (
	"net/http"

	"github.com/ashez2000/rssaggr/internal/database"
)

type Application struct {
	DB *database.Queries
}

func (app *Application) hello(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, "Hello, world!")
}
