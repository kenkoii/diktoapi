package routers

import (
	"github.com/gorilla/mux"
	"github.com/kenkoii/diktoapi/api/handlers"
)

// SetWordsRoutes sets routing for Words Endpoint
func SetWordsRoutes(router *mux.Router) *mux.Router {
	wordsRouter := mux.NewRouter()
	wordsRouter.HandleFunc("/api/v1/words", handlers.PostWordEndpoint).Methods("POST")
	wordsRouter.HandleFunc("/api/v1/words/{id}", handlers.GetWordEndpoint).Methods("GET")
	wordsRouter.HandleFunc("/api/v1/words/{id}", handlers.UpdateWordEndpoint).Methods("PUT")
	wordsRouter.HandleFunc("/api/v1/words/{id}/favorite/{userid}", handlers.FavoriteWordEndpoint).Methods("GET")
	wordsRouter.HandleFunc("/api/v1/words/{id}/favorite/{userid}", handlers.RemoveFavoriteWordEndpoint).Methods("DELETE")
	wordsRouter.HandleFunc("/api/v1/words/{id}/lemma", handlers.GetLemmaEndpoint).Methods("GET")
	router.PathPrefix("/api/v1/words").Handler(wordsRouter)
	//router.
	return router
}
