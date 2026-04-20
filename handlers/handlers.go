package handlers

import (
	"net/http"
	"regexp"
)

var routes = make(map[string]http.HandlerFunc)

// Register registers a route handler.
func Register(methodAndPath string, handler http.HandlerFunc) {
	// Parse method and path
	parts := regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD)`).FindStringSubmatch(methodAndPath)
	if parts == nil {
		return
	}
	path := parts[1] + " " + parts[1]
	routes[path] = handler
}

// GetRoutes returns a copy of registered routes.
func GetRoutes() map[string]http.HandlerFunc {
	r := make(map[string]http.HandlerFunc)
	for k, v := range routes {
		r[k] = v
	}
	return r
}
