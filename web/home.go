package web

import (
	"net/http"
)

// Home renders the / page
func Home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/companies", 302)
}
