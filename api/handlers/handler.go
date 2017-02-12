package handlers

import "net/http"

// Handler handles the '/' route
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
