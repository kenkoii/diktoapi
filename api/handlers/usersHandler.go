package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kenkoii/diktoapi/api/models"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// GetUserEndpoint handles the /api/v1/users/{id} {GET} method
func GetUserEndpoint(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	password, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := appengine.NewContext(c.Request)
	user, err := models.GetUser(ctx, id, password)
	if err != nil {
		LogErrorGin(c, err)
	}
	c.JSON(200, user)
}

// func GetUserEndpoint(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, _ := strconv.ParseInt(vars["id"], 10, 64)
// 	password, _ := strconv.ParseInt(vars["password"], 10, 64)
// 	ctx := appengine.NewContext(r)
// 	user, err := models.GetUser(ctx, id, password)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(user)
// }

// UpdateUserEndpoint handles the /api/v1/users/{id} {GET} method
func UpdateUserEndpoint(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	ctx := appengine.NewContext(c.Request)
	user, err := models.UpdateUser(ctx, id, c.Request.Body)
	if err != nil {
		LogErrorGin(c, err)
	}
	c.JSON(http.StatusOK, user)
}

// func UpdateUserEndpoint(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, _ := strconv.ParseInt(vars["id"], 10, 64)

// 	ctx := appengine.NewContext(r)
// 	user, err := models.UpdateUser(ctx, id, r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(user)
// }

func LogErrorGin(c *gin.Context, err error) {
	log.Errorf(c, "could not put into datastore: %v", err)
	c.String(http.StatusOK, "-1")
}
