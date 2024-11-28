package main

import (
	"log"
	"net/http"

	"github.com/ashez2000/rssaggr/internal/database"
)

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
