package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkoii/diktoapi/api/handlers"
)

func InitGinRoutes(router *gin.Engine) *gin.Engine {
	// v1 := router.Group("/api/v1/analytics")
	v1 := router.Group("/api/v1")
	v1.POST("/words", handlers.PostWordEndpoint)
	v1.PUT("/words/:id", handlers.UpdateWordEndpoint)
	v1.GET("/words/:id", handlers.GetWordEndpoint)
	v1.GET("/words/:id/lemma", handlers.GetLemmaEndpoint)
	v1.POST("/favorite", handlers.FavoriteWordEndpoint)
	v1.POST("/favorite/remove", handlers.RemoveFavoriteWordEndpoint)
	v1.POST("/favorite/frontend", handlers.FrontendFavoriteWordEndpoint)
	v1.POST("/favorite/frontend/remove", handlers.RemoveFavoriteWordEndpoint)
	v1.GET("/users/:id/:password", handlers.GetUserEndpoint)
	v1.PUT("/users/:id", handlers.UpdateUserEndpoint)

	//FRONTEND
	public := router.Group("/")
	public.GET("/list/:id/:password", handlers.ListHandler)
	public.GET("/word/:word", handlers.DetailHandler)
	public.GET("/word/:word/:id/:password", handlers.DetailHandler)

	return router
}
