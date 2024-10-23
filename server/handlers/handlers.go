// /pkg/handlers/handlers.go

package handlers

import (
	"net/http"
	"rtforum/server/sessions"
)

// HandleSPA serves the main HTML file for all non-API routes
func HandleSPA(w http.ResponseWriter, r *http.Request) {
	_, valid := sessions.ValidateSession(r)
	if !valid {
		// If the user is not authenticated, return 401 Unauthorized for testing purposes
		if r.Header.Get("Test-Mode") == "true" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Otherwise, redirect to the login page (for production usage)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "static/index.html")

}
