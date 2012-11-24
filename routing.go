package webshell

import "net/http"

type RouteHandler func(w http.ResponseWriter, r *http.Request)

var RouterMux = http.NewServeMux()

// AddRoute is syntactic sugar for adding routes to aid in late night hacks.
// All this does is call RouterMux.HandleFunc(path, handler).
func AddRoute(path string, handler RouteHandler) {
	RouterMux.HandleFunc(path, handler)
}

// StaticRoute sets up a route for serving static files.
// route sets the route that should be used, and path is the path to the
// static files
func StaticRoute(route string, path string) {
        RouterMux.Handle(route, http.FileServer(http.Dir(path)))
}
