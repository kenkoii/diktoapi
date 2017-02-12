package routers

import (
	"github.com/gorilla/mux"
)

// InitRoutes is for initializing all routes/endpoints
func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetCommentsRoutes(router)
	return router
}
