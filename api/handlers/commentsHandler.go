package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kenkoii/diktoapi/api/models"
	"google.golang.org/appengine"
)

func PostCommentEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := appengine.NewContext(r)
	comment, err := models.NewComment(ctx, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(comment)
}
