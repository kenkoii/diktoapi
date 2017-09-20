package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kenkoii/diktoapi/api/models"
	"google.golang.org/appengine"
)

// Handler handles the '/' route
func ListHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	password, _ := strconv.ParseInt(c.Param("password"), 10, 64)
	ctx := appengine.NewContext(c.Request)
	user, err := models.GetUser(ctx, id, password)
	if err != nil {
		LogErrorGin(c, err)
	}
	RenderResult(c, user, "list.tmpl", gin.H{})
}

func SettingsHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	password, _ := strconv.ParseInt(c.Param("password"), 10, 64)
	ctx := appengine.NewContext(c.Request)
	user, err := models.GetUser(ctx, id, password)
	if err != nil {
		LogErrorGin(c, err)
	}
	RenderResult(c, user, "settings.tmpl", gin.H{})
}

func DetailHandler(c *gin.Context) {
	// id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	// password, _ := strconv.ParseInt(c.Param("password"), 10, 64)
	var user interface{}
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	password, _ := strconv.ParseInt(c.Query("password"), 10, 64)
	word := c.Param("word")
	ctx := appengine.NewContext(c.Request)
	w, err := models.GetWord(ctx, word)
	if err != nil {
		LogErrorGin(c, err)
	}
	if id != 0 && password != 0 {
		user, err = models.GetUser(ctx, id, password)
		if err != nil {
			LogErrorGin(c, err)
		}
	}
	RenderResult(c, gin.H{"word": w, "user": user}, "detail.tmpl", gin.H{})
}

func RenderResult(c *gin.Context, results interface{}, template string, params gin.H) {
	PrintObject(results)
	c.HTML(
		http.StatusOK,
		template,
		gin.H{
			"title":  "Main website",
			"method": "POST",
			"data":   results,
			"params": params,
		})
}

func PrintObject(obj interface{}) {
	log.Print("===Object===")
	log.Println(obj)
	log.Print("============")
}
