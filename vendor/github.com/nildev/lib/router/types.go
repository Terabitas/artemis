package router

import "net/http"

type (
	// Route structure contains information about a route and a reference to a handler.
	Route struct {
		Name        string
		Method      []string
		Pattern     string
		HandlerFunc http.HandlerFunc
		Queries     []string
		Protected   bool
	}

	// Routes is a collection of Route structures.
	Routes struct {
		BasePattern string
		Routes      []Route
	}
)
