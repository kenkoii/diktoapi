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
	// gin.SetMode(gin.ReleaseMode)
	router.Use(middlewares.CORSMiddleware())
	router.Static("assets", "assets")
	router.LoadHTMLGlob("templates/*")
	// router.LoadHTMLGlob("github.com/kenkoii/diktoapi/api/templates/*")
	// router.LoadHTMLGlob(filepath.Join(os.Getenv("GOPATH"), "src/github.com/kenkoii/diktoapi/api/templates/*"))
	router = routers.InitGinRoutes(router)
	return router
}
