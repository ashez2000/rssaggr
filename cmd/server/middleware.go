package main

import (
	"net/http"

	"github.com/ashez2000/rssaggr/internal/auth"
	"github.com/ashez2000/rssaggr/internal/database"
)

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
