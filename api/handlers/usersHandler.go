package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kenkoii/diktoapi/api/models"
	"google.golang.org/appengine"
)

// GetUserEndpoint handles the /api/v1/users/{id} {GET} method
func GetUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	ctx := appengine.NewContext(r)
	user, err := models.GetUser(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// UpdateUserEndpoint handles the /api/v1/users/{id} {GET} method
func UpdateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	ctx := appengine.NewContext(r)
	user, err := models.UpdateUser(ctx, id, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}
