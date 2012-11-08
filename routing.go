package webshell

import "net/http"

type RouteHandler func (w http.ResponseWriter, r *http.Request)

var RouterMux = http.NewServeMux()

// AddRoute is syntactic sugar for adding routes to aid in late night hacks.
// All this does is call RouterMux.HandleFunc(path, handler).
func AddRoute(path string, handler RouteHandler) {
        RouterMux.HandleFunc(path, handler)
}
