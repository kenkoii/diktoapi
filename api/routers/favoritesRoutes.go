package routers

import (
	"github.com/gorilla/mux"
	"github.com/kenkoii/diktoapi/api/handlers"
)

// SetFavoritesRoutes sets routing for Words Endpoint
func SetFavoritesRoutes(router *mux.Router) *mux.Router {
	favoritesRouter := mux.NewRouter()
	favoritesRouter.HandleFunc("/api/v1/favorite", handlers.FavoriteWordEndpoint).Methods("POST")
	favoritesRouter.HandleFunc("/api/v1/favorite/remove", handlers.RemoveFavoriteWordEndpoint).Methods("POST")
	router.PathPrefix("/api/v1/favorite").Handler(favoritesRouter)
	//router.
	return router
}