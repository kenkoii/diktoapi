package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenkoii/diktoapi/api/middlewares"
	"github.com/kenkoii/diktoapi/api/routers"
	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", GetMainEngine())

	/*
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
			log.Printf("Defaulting to port %s", port)
		}

		log.Printf("Listening on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	*/
	appengine.Main()
}

func GetMainEngine() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	router.Static("assets", "assets")
	router.LoadHTMLGlob("templates/*")
	router = routers.InitGinRoutes(router)
	return router
}
