package router

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a gorilla/mux router from a list of Route structures.
func NewRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes.Routes {
		router.
			Methods(route.Method...).
			Name(route.Name).
			Path(routes.BasePattern + route.Pattern).
			Queries(route.Queries...).
			HandlerFunc(route.HandlerFunc)
	}

	return router
}
