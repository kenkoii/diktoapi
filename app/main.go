package app

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/kenkoii/diktoapi/api/handlers"
	"github.com/kenkoii/diktoapi/api/routers"
	"github.com/rs/cors"
)

func init() {
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
	})
	router := routers.InitRoutes()
	router.HandleFunc("/", handlers.Handler)
	router.HandleFunc("/worker", handlers.Worker)
	n := negroni.Classic()
	handler := c.Handler(router)
	n.UseHandler(handler)
	http.Handle("/", n)
}
