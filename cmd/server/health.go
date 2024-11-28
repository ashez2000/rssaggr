package main

import "net/http"

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, "OK")
}
