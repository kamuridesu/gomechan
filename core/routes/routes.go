package routes

import "net/http"

// Interface for routers like http and/or http.ServeMux
// This allow us to use common methods for those routers
type router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Add a health check route to the router
// It can be consulted by requesting the /health path
func AddHealthCheck(r router) {
	r.HandleFunc("/health", HealthCheck)
}
