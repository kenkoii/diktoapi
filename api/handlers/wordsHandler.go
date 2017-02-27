package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kenkoii/diktoapi/api/models"
	"google.golang.org/appengine"
)

// PostWordEndpoint handles the /api/v1/words/ {POST} method
func PostWordEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := appengine.NewContext(r)
	word, err := models.NewWord(ctx, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(word)
}

// GetWordEndpoint handles the /api/v1/words/{id} {GET} method
func GetWordEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ctx := appengine.NewContext(r)
	word, err := models.GetWord(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(word)
}
