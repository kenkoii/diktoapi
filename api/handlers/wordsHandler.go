package handlers

import (
	"encoding/json"
	"fmt"
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

// UpdateWordEndpoint handles the /api/v1/words/ {POST} method
func UpdateWordEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := appengine.NewContext(r)
	word, err := models.UpdateWord(ctx, r.Body)
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

// GetLemmaEndpoint handles the /api/v1/words/{id} {GET} method
func GetLemmaEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ctx := appengine.NewContext(r)
	word, err := models.GetLemma(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, word.Lemma)
}

// FavoriteWordEndpoint handles the /api/v1/words/{id} {POST} method
func FavoriteWordEndpoint(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]
	// userid, _ := strconv.ParseInt(vars["userid"], 10, 64)
	ctx := appengine.NewContext(r)
	res, err := models.FavoriteWord(ctx, r.Body)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintln(w, res)
	// json.NewEncoder(w).Encode(res)
}

// RemoveFavoriteWordEndpoint handles the /api/v1/words/{id} {POST} method
func RemoveFavoriteWordEndpoint(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]
	// userid, _ := strconv.ParseInt(vars["userid"], 10, 64)
	ctx := appengine.NewContext(r)
	word, err := models.RemoveFavoriteWord(ctx, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(word)
}

// FrontendFavoriteWordEndpoint handles the /api/v1/words/{id} {POST} method
func FrontendFavoriteWordEndpoint(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]
	// userid, _ := strconv.ParseInt(vars["userid"], 10, 64)
	ctx := appengine.NewContext(r)
	res, err := models.FrontendFavoriteWord(ctx, r.Body)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}
	json.NewEncoder(w).Encode(res)
}
