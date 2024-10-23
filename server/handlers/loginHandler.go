package handlers

import (
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")

}
