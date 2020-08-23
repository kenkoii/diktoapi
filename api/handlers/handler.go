package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kenkoii/diktoapi/api/models"
)

// Handler handles the '/' route
func ListHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	password, _ := strconv.ParseInt(c.Param("password"), 10, 64)
	ctx := c.Request.Context()
	user, err := models.GetUser(ctx, id, password)
	if err != nil {
		// LogErrorGin(c, err)
		RenderResult(c, gin.H{"error": err.Error(), "code": http.StatusForbidden}, "error.tmpl", http.StatusForbidden)
		return
	}
	RenderResult(c, user, "list.tmpl", http.StatusOK)
}

func SettingsHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	password, _ := strconv.ParseInt(c.Param("password"), 10, 64)
	ctx := c.Request.Context()
	user, err := models.GetUser(ctx, id, password)
	if err != nil {
		LogErrorGin(c, err)
	}
	RenderResult(c, user, "settings.tmpl", http.StatusOK)
}

func DetailHandler(c *gin.Context) {
	// id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	// password, _ := strconv.ParseInt(c.Param("password"), 10, 64)
	var user interface{}
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	password, _ := strconv.ParseInt(c.Query("password"), 10, 64)
	word := c.Param("word")
	ctx := c.Request.Context()
	w, err := models.GetWord(ctx, word)
	if err != nil {
		// LogErrorGin(c, err)
		RenderResult(c, gin.H{"error": "Word not found", "code": http.StatusNotFound}, "error.tmpl", http.StatusNotFound)
		return
	}
	if id != 0 && password != 0 {
		user, err = models.GetUser(ctx, id, password)
		PrintObject(user)
		if err != nil {
			// LogErrorGin(c, err)
			RenderResult(c, gin.H{"error": err.Error(), "code": http.StatusForbidden}, "error.tmpl", http.StatusForbidden)
			return
		}
	}
	RenderResult(c, gin.H{"word": w, "user": user}, "detail.tmpl", http.StatusOK)
}

func RenderResult(c *gin.Context, results interface{}, template string, statusCode int) {
	PrintObject(results)
	c.HTML(
		statusCode,
		template,
		gin.H{
			"title":  "Main website",
			"method": "POST",
			"data":   results,
		})
}

func PrintObject(obj interface{}) {
	log.Print("===Object===")
	log.Println(obj)
	log.Print("============")
}
