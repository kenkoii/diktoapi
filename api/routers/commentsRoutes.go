package routers

import (
	"github.com/gorilla/mux"
	"github.com/kenkoii/diktoapi/api/handlers"
)

func SetCommentsRoutes(router *mux.Router) *mux.Router {
	commentsRouter := mux.NewRouter()
	commentsRouter.HandleFunc("/api/v1/comments", handlers.PostCommentEndpoint).Methods("POST")
	router.PathPrefix("/api/v1/comments").Handler(commentsRouter)
	//router.
	return router
}
