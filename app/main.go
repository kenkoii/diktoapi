package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenkoii/diktoapi/api/middlewares"
	"github.com/kenkoii/diktoapi/api/routers"
)

func init() {
	http.Handle("/", GetMainEngine())
}

func GetMainEngine() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	router.Static("assets", "../api/assets")
	router.LoadHTMLGlob("../api/templates/*")
	router = routers.InitGinRoutes(router)
	return router
}
